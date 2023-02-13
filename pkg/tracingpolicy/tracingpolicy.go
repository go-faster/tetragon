// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package tracingpolicy

import (
	"github.com/cilium/tetragon/pkg/k8s/apis/cilium.io/v1alpha1"
)

// revive:disable

// TracingPolicy is an interface for a tracing policy
// This is implemented by v1alpha1.types.TracingPolicy and
// config.GenericTracingConf. The former is what is the k8s API server uses,
// and the latter is used when we load files directly (e.g., via the cli).
type TracingPolicy interface {
	// TpName returns the name of the policy.
	TpName() string
	// TpSpec  returns the specification of the policy
	TpSpec() *v1alpha1.TracingPolicySpec
	// TpInfo returns a description of the policy
	TpInfo() string
}

// TracingPolicyNamespaced is a tracing policy applied on a specific namespace
type TracingPolicyNamespaced interface {
	TracingPolicy
	// TpNamespace returns the namespace of the policy
	TpNamespace() string
}

type SpecTp struct {
	name string
	info string
	spec *v1alpha1.TracingPolicySpec
}

func NewTp(name string, info string, spec *v1alpha1.TracingPolicySpec) *SpecTp {
	return &SpecTp{
		name: name,
		info: info,
		spec: spec,
	}
}

func (s *SpecTp) TpName() string {
	return s.name
}

func (s *SpecTp) TpInfo() string {
	return s.info
}

func (s *SpecTp) TpSpec() *v1alpha1.TracingPolicySpec {
	return s.spec
}
