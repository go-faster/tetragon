// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package opcodemetrics

import (
	"fmt"

	"github.com/go-faster/tetragon/pkg/metrics/consts"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	MsgOpsCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        consts.MetricNamePrefix + "msg_op_total",
		Help:        "The total number of times we encounter a given message opcode. For internal use only.",
		ConstLabels: nil,
	}, []string{"msg_op"})
)

// Get a new handle on a msgOpsCount metric for an OpCode
func GetOpTotal(op int) prometheus.Counter {
	return MsgOpsCount.WithLabelValues(fmt.Sprint(op))
}

// Increment an msgOpsCount for an OpCode
func OpTotalInc(op int) {
	GetOpTotal(op).Inc()
}
