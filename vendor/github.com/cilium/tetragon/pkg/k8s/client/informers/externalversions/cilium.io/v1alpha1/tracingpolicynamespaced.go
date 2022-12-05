// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	ciliumiov1alpha1 "github.com/cilium/tetragon/pkg/k8s/apis/cilium.io/v1alpha1"
	versioned "github.com/cilium/tetragon/pkg/k8s/client/clientset/versioned"
	internalinterfaces "github.com/cilium/tetragon/pkg/k8s/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/cilium/tetragon/pkg/k8s/client/listers/cilium.io/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// TracingPolicyNamespacedInformer provides access to a shared informer and lister for
// TracingPoliciesNamespaced.
type TracingPolicyNamespacedInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.TracingPolicyNamespacedLister
}

type tracingPolicyNamespacedInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewTracingPolicyNamespacedInformer constructs a new informer for TracingPolicyNamespaced type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewTracingPolicyNamespacedInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredTracingPolicyNamespacedInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredTracingPolicyNamespacedInformer constructs a new informer for TracingPolicyNamespaced type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredTracingPolicyNamespacedInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CiliumV1alpha1().TracingPoliciesNamespaced(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CiliumV1alpha1().TracingPoliciesNamespaced(namespace).Watch(context.TODO(), options)
			},
		},
		&ciliumiov1alpha1.TracingPolicyNamespaced{},
		resyncPeriod,
		indexers,
	)
}

func (f *tracingPolicyNamespacedInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredTracingPolicyNamespacedInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *tracingPolicyNamespacedInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&ciliumiov1alpha1.TracingPolicyNamespaced{}, f.defaultInformer)
}

func (f *tracingPolicyNamespacedInformer) Lister() v1alpha1.TracingPolicyNamespacedLister {
	return v1alpha1.NewTracingPolicyNamespacedLister(f.Informer().GetIndexer())
}
