// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package client

import (
	"context"
	_ "embed"
	goerrors "errors"
	"fmt"
	"time"

	k8sversion "github.com/go-faster/tetragon/pkg/k8s/version"

	"github.com/cilium/cilium/pkg/logging"
	"github.com/cilium/cilium/pkg/logging/logfields"
	"github.com/cilium/cilium/pkg/versioncheck"
	ciliumio "github.com/go-faster/tetragon/pkg/k8s/apis/cilium.io"
	"github.com/go-faster/tetragon/pkg/k8s/apis/cilium.io/v1alpha1"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/yaml"
)

const (
	// subsysK8s is the value for logfields.LogSubsys
	subsysK8s = "k8s"

	// CustomResourceDefinitionSchemaVersionKey is key to label which holds the CRD schema version
	CustomResourceDefinitionSchemaVersionKey = ciliumio.GroupName + ".k8s.crd.schema.version"
)

var (
	// log is the k8s package logger object.
	log = logging.DefaultLogger.WithField(logfields.LogSubsys, subsysK8s)

	comparableCRDSchemaVersion = versioncheck.MustVersion(v1alpha1.CustomResourceDefinitionSchemaVersion)
)

// CreateCustomResourceDefinitions creates our CRD objects in the Kubernetes
// cluster.
func CreateCustomResourceDefinitions(clientset apiextensionsclient.Interface) error {
	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return createTPCRDs(clientset)
	})

	return g.Wait()
}

var (
	//go:embed crds/v1alpha1/cilium.io_tracingpolicies.yaml
	crdsv1Alpha1TracingPolicies []byte

	//go:embed crds/v1alpha1/cilium.io_tracingpoliciesnamespaced.yaml
	crdsv1Alpha1TracingPoliciesNamespaced []byte
)

// GetPregeneratedCRD returns the pregenerated CRD based on the requested CRD
// name. The pregenerated CRDs are generated by the controller-gen tool and
// serialized into binary form by go-bindata. This function retrieves CRDs from
// the binary form.
func GetPregeneratedCRD(crdName string) apiextensionsv1.CustomResourceDefinition {
	var (
		err      error
		crdBytes []byte
	)

	scopedLog := log.WithField("crdName", crdName)

	switch crdName {
	case v1alpha1.TPCRDName:
		crdBytes = crdsv1Alpha1TracingPolicies
	case v1alpha1.TPNamespacedCRDName:
		crdBytes = crdsv1Alpha1TracingPoliciesNamespaced
	default:
		scopedLog.Fatal("Pregenerated CRD does not exist")
	}

	isoCRD := apiextensionsv1.CustomResourceDefinition{}
	err = yaml.Unmarshal(crdBytes, &isoCRD)
	if err != nil {
		scopedLog.WithError(err).Fatal("Error unmarshalling pregenerated CRD")
	}

	return isoCRD
}

func createTPCRDs(clientset apiextensionsclient.Interface) error {
	// custer-wide tracing policy CRD
	isoCRD := GetPregeneratedCRD(v1alpha1.TPCRDName)
	if err := createUpdateCRD(
		clientset,
		v1alpha1.TPCRDName,
		constructV1CRD(v1alpha1.TPName, isoCRD),
		newDefaultPoller(),
	); err != nil {
		return err
	}

	// namespaced tracing policy CRD
	isoCRD = GetPregeneratedCRD(v1alpha1.TPNamespacedCRDName)
	if err := createUpdateCRD(
		clientset,
		v1alpha1.TPNamespacedCRDName,
		constructV1CRD(v1alpha1.TPNamespacedName, isoCRD),
		newDefaultPoller(),
	); err != nil {
		return err
	}
	return nil
}

// createUpdateCRD ensures the CRD object is installed into the K8s cluster. It
// will create or update the CRD and its validation schema as necessary. This
// function only accepts v1 CRD objects, and defers to its v1beta1 variant if
// the cluster only supports v1beta1 CRDs. This allows us to convert all our
// CRDs into v1 form and only perform conversions on-demand, simplifying the
// code.
func createUpdateCRD(
	clientset apiextensionsclient.Interface,
	crdName string,
	crd *apiextensionsv1.CustomResourceDefinition,
	poller poller,
) error {
	scopedLog := log.WithField("name", crdName)

	if !k8sversion.Capabilities().APIExtensionsV1CRD {
		log.Infof("K8s apiserver does not support v1 CRDs, falling back to v1beta1")

		return createUpdateV1beta1CRD(
			scopedLog,
			clientset.ApiextensionsV1beta1(),
			crdName,
			crd,
			poller,
		)
	}

	v1CRDClient := clientset.ApiextensionsV1()
	clusterCRD, err := v1CRDClient.CustomResourceDefinitions().Get(
		context.TODO(),
		crd.ObjectMeta.Name,
		metav1.GetOptions{})
	if errors.IsNotFound(err) {
		scopedLog.Info("Creating CRD (CustomResourceDefinition)...")

		clusterCRD, err = v1CRDClient.CustomResourceDefinitions().Create(
			context.TODO(),
			crd,
			metav1.CreateOptions{})
		// This occurs when multiple agents race to create the CRD. Since another has
		// created it, it will also update it, hence the non-error return.
		if errors.IsAlreadyExists(err) {
			return nil
		}
	}
	if err != nil {
		return err
	}

	if err := updateV1CRD(scopedLog, crd, clusterCRD, v1CRDClient, poller); err != nil {
		return err
	}
	if err := waitForV1CRD(scopedLog, crdName, clusterCRD, v1CRDClient, poller); err != nil {
		return err
	}

	scopedLog.Info("CRD (CustomResourceDefinition) is installed and up-to-date")

	return nil
}

func constructV1CRD(
	name string,
	template apiextensionsv1.CustomResourceDefinition,
) *apiextensionsv1.CustomResourceDefinition {
	return &apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				CustomResourceDefinitionSchemaVersionKey: v1alpha1.CustomResourceDefinitionSchemaVersion,
			},
		},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group: ciliumio.GroupName,
			Names: apiextensionsv1.CustomResourceDefinitionNames{
				Kind:       template.Spec.Names.Kind,
				Plural:     template.Spec.Names.Plural,
				ShortNames: template.Spec.Names.ShortNames,
				Singular:   template.Spec.Names.Singular,
			},
			Scope:    template.Spec.Scope,
			Versions: template.Spec.Versions,
		},
	}
}

func needsUpdateV1(clusterCRD *apiextensionsv1.CustomResourceDefinition) bool {
	if clusterCRD.Spec.Versions[0].Schema == nil {
		// no validation detected
		return true
	}
	v, ok := clusterCRD.Labels[CustomResourceDefinitionSchemaVersionKey]
	if !ok {
		// no schema version detected
		return true
	}

	clusterVersion, err := versioncheck.Version(v)
	if err != nil || clusterVersion.LT(comparableCRDSchemaVersion) {
		// version in cluster is either unparsable or smaller than current version
		return true
	}

	return false
}

func updateV1CRD(
	scopedLog *logrus.Entry,
	crd, clusterCRD *apiextensionsv1.CustomResourceDefinition,
	client v1client.CustomResourceDefinitionsGetter,
	poller poller,
) error {
	scopedLog.Debug("Checking if CRD (CustomResourceDefinition) needs update...")

	if crd.Spec.Versions[0].Schema != nil && needsUpdateV1(clusterCRD) {
		scopedLog.Info("Updating CRD (CustomResourceDefinition)...")

		// Update the CRD with the validation schema.
		err := poller.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
			var err error
			clusterCRD, err = client.CustomResourceDefinitions().Get(
				context.TODO(),
				crd.ObjectMeta.Name,
				metav1.GetOptions{})
			if err != nil {
				return false, err
			}

			// This seems too permissive but we only get here if the version is
			// different per needsUpdate above. If so, we want to update on any
			// validation change including adding or removing validation.
			if needsUpdateV1(clusterCRD) {
				scopedLog.Debug("CRD validation is different, updating it...")

				clusterCRD.ObjectMeta.Labels = crd.ObjectMeta.Labels
				clusterCRD.Spec = crd.Spec

				// Even though v1 CRDs omit this field by default (which also
				// means it's false) it is still carried over from the previous
				// CRD. Therefore, we must set this to false explicitly because
				// the apiserver will carry over the old value (true).
				clusterCRD.Spec.PreserveUnknownFields = false

				_, err := client.CustomResourceDefinitions().Update(
					context.TODO(),
					clusterCRD,
					metav1.UpdateOptions{})
				switch {
				case errors.IsConflict(err): // Occurs as Operators race to update CRDs.
					scopedLog.WithError(err).
						Debug("The CRD update was based on an older version, retrying...")
					return false, nil
				case err == nil:
					return true, nil
				}

				scopedLog.WithError(err).Debug("Unable to update CRD validation")

				return false, err
			}

			return true, nil
		})
		if err != nil {
			scopedLog.WithError(err).Error("Unable to update CRD")
			return err
		}
	}

	return nil
}

func waitForV1CRD(
	scopedLog *logrus.Entry,
	crdName string,
	crd *apiextensionsv1.CustomResourceDefinition,
	client v1client.CustomResourceDefinitionsGetter,
	poller poller,
) error {
	scopedLog.Debug("Waiting for CRD (CustomResourceDefinition) to be available...")

	err := poller.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
		for _, cond := range crd.Status.Conditions {
			switch cond.Type {
			case apiextensionsv1.Established:
				if cond.Status == apiextensionsv1.ConditionTrue {
					return true, nil
				}
			case apiextensionsv1.NamesAccepted:
				if cond.Status == apiextensionsv1.ConditionFalse {
					err := goerrors.New(cond.Reason)
					scopedLog.WithError(err).Error("Name conflict for CRD")
					return false, err
				}
			}
		}

		var err error
		if crd, err = client.CustomResourceDefinitions().Get(
			context.TODO(),
			crd.ObjectMeta.Name,
			metav1.GetOptions{}); err != nil {
			return false, err
		}
		return false, err
	})
	if err != nil {
		return fmt.Errorf("error occurred waiting for CRD: %w", err)
	}

	return nil
}

// poller is an interface that abstracts the polling logic when dealing with
// CRD changes / updates to the apiserver. The reason this exists is mainly for
// unit-testing.
type poller interface {
	Poll(interval, duration time.Duration, conditionFn func() (bool, error)) error
}

func newDefaultPoller() defaultPoll {
	return defaultPoll{}
}

type defaultPoll struct{}

func (p defaultPoll) Poll(
	interval, duration time.Duration,
	conditionFn func() (bool, error),
) error {
	return wait.Poll(interval, duration, conditionFn)
}

// RegisterCRDs registers all CRDs with the K8s apiserver.
func RegisterCRDs(clientset apiextensionsclient.Interface) error {
	if err := CreateCustomResourceDefinitions(clientset); err != nil {
		return fmt.Errorf("Unable to create custom resource definition: %w", err)
	}

	return nil
}
