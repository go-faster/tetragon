// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/go-faster/tetragon/pkg/k8s/client/clientset/versioned/typed/cilium.io/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeCiliumV1alpha1 struct {
	*testing.Fake
}

func (c *FakeCiliumV1alpha1) PodInfo(namespace string) v1alpha1.PodInfoInterface {
	return newFakePodInfo(c, namespace)
}

func (c *FakeCiliumV1alpha1) TracingPolicies() v1alpha1.TracingPolicyInterface {
	return newFakeTracingPolicies(c)
}

func (c *FakeCiliumV1alpha1) TracingPoliciesNamespaced(namespace string) v1alpha1.TracingPolicyNamespacedInterface {
	return newFakeTracingPoliciesNamespaced(c, namespace)
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeCiliumV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
