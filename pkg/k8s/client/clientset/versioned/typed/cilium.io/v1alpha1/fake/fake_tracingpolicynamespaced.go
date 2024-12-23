// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/go-faster/tetragon/pkg/k8s/apis/cilium.io/v1alpha1"
	ciliumiov1alpha1 "github.com/go-faster/tetragon/pkg/k8s/client/clientset/versioned/typed/cilium.io/v1alpha1"
	gentype "k8s.io/client-go/gentype"
)

// fakeTracingPoliciesNamespaced implements TracingPolicyNamespacedInterface
type fakeTracingPoliciesNamespaced struct {
	*gentype.FakeClientWithList[*v1alpha1.TracingPolicyNamespaced, *v1alpha1.TracingPolicyNamespacedList]
	Fake *FakeCiliumV1alpha1
}

func newFakeTracingPoliciesNamespaced(fake *FakeCiliumV1alpha1, namespace string) ciliumiov1alpha1.TracingPolicyNamespacedInterface {
	return &fakeTracingPoliciesNamespaced{
		gentype.NewFakeClientWithList[*v1alpha1.TracingPolicyNamespaced, *v1alpha1.TracingPolicyNamespacedList](
			fake.Fake,
			namespace,
			v1alpha1.SchemeGroupVersion.WithResource("tracingpoliciesnamespaced"),
			v1alpha1.SchemeGroupVersion.WithKind("TracingPolicyNamespaced"),
			func() *v1alpha1.TracingPolicyNamespaced { return &v1alpha1.TracingPolicyNamespaced{} },
			func() *v1alpha1.TracingPolicyNamespacedList { return &v1alpha1.TracingPolicyNamespacedList{} },
			func(dst, src *v1alpha1.TracingPolicyNamespacedList) { dst.ListMeta = src.ListMeta },
			func(list *v1alpha1.TracingPolicyNamespacedList) []*v1alpha1.TracingPolicyNamespaced {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v1alpha1.TracingPolicyNamespacedList, items []*v1alpha1.TracingPolicyNamespaced) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
