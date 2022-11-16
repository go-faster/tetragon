// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package rthooks

import (
	"context"
	"testing"

	"github.com/cilium/tetragon/api/v1/tetragon"
)

type DummyHookRunner struct {
	*testing.T
}

func (o DummyHookRunner) RunHooks(ctx context.Context, req *tetragon.RuntimeHookRequest) error {
	return nil
}
