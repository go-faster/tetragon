// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	ciliumiov1alpha1 "github.com/go-faster/tetragon/pkg/k8s/apis/cilium.io/v1alpha1"
	versioned "github.com/go-faster/tetragon/pkg/k8s/client/clientset/versioned"
	internalinterfaces "github.com/go-faster/tetragon/pkg/k8s/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/go-faster/tetragon/pkg/k8s/client/listers/cilium.io/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// TracingPolicyInformer provides access to a shared informer and lister for
// TracingPolicies.
type TracingPolicyInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.TracingPolicyLister
}

type tracingPolicyInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewTracingPolicyInformer constructs a new informer for TracingPolicy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewTracingPolicyInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredTracingPolicyInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredTracingPolicyInformer constructs a new informer for TracingPolicy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredTracingPolicyInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CiliumV1alpha1().TracingPolicies().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CiliumV1alpha1().TracingPolicies().Watch(context.TODO(), options)
			},
		},
		&ciliumiov1alpha1.TracingPolicy{},
		resyncPeriod,
		indexers,
	)
}

func (f *tracingPolicyInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredTracingPolicyInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *tracingPolicyInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&ciliumiov1alpha1.TracingPolicy{}, f.defaultInformer)
}

func (f *tracingPolicyInformer) Lister() v1alpha1.TracingPolicyLister {
	return v1alpha1.NewTracingPolicyLister(f.Informer().GetIndexer())
}
