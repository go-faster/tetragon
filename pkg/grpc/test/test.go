// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon
package test

import (
	"github.com/go-faster/tetragon/api/v1/tetragon"
	"github.com/go-faster/tetragon/pkg/api/ops"
	"github.com/go-faster/tetragon/pkg/api/testapi"
	"github.com/go-faster/tetragon/pkg/eventcache"
	"github.com/go-faster/tetragon/pkg/ktime"
	"github.com/go-faster/tetragon/pkg/logger"
	"github.com/go-faster/tetragon/pkg/process"
	"github.com/go-faster/tetragon/pkg/reader/node"
	"github.com/go-faster/tetragon/pkg/reader/notify"
)

var (
	nodeName = node.GetNodeNameForExport()
)

type MsgTestEventUnix struct {
	testapi.MsgTestEvent
}

func (msg *MsgTestEventUnix) Notify() bool {
	return true
}

func (msg *MsgTestEventUnix) RetryInternal(ev notify.Event, timestamp uint64) (*process.ProcessInternal, error) {
	return eventcache.HandleGenericInternal(ev, 0, timestamp)
}

func (msg *MsgTestEventUnix) Retry(internal *process.ProcessInternal, ev notify.Event) error {
	return eventcache.HandleGenericEvent(internal, ev)
}

func (msg *MsgTestEventUnix) HandleMessage() *tetragon.GetEventsResponse {
	var res *tetragon.GetEventsResponse
	switch msg.Common.Op {
	case ops.MSG_OP_TEST:
		res = &tetragon.GetEventsResponse{
			Event: &tetragon.GetEventsResponse_Test{Test: &tetragon.Test{
				Arg0: msg.Arg0,
				Arg1: msg.Arg1,
				Arg2: msg.Arg2,
				Arg3: msg.Arg3,
			}},
			NodeName: nodeName,
			Time:     ktime.ToProto(msg.Common.Ktime),
		}
	default:
		logger.GetLogger().WithField("message", msg).Warn("HandleTestMessage: Unhandled event")
	}
	return res
}

func (msg *MsgTestEventUnix) Cast(o interface{}) notify.Message {
	t := o.(testapi.MsgTestEvent)
	return &MsgTestEventUnix{MsgTestEvent: t}
}
