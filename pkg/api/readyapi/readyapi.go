// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon
package readyapi

import (
	"fmt"

	"github.com/go-faster/tetragon/api/v1/tetragon"
	"github.com/go-faster/tetragon/pkg/process"
	"github.com/go-faster/tetragon/pkg/reader/notify"
)

type MsgTetragonReady struct{}

func (msg *MsgTetragonReady) Notify() bool {
	return false
}

func (msg *MsgTetragonReady) RetryInternal(ev notify.Event, timestamp uint64) (*process.ProcessInternal, error) {
	return nil, fmt.Errorf("Unsupported cache event MsgTetragonReady")
}

func (msg *MsgTetragonReady) Retry(internal *process.ProcessInternal, ev notify.Event) error {
	return fmt.Errorf("Unsupported cache retry event MsgTetragonReady")
}

func (msg *MsgTetragonReady) HandleMessage() *tetragon.GetEventsResponse {
	return nil
}

func (msg *MsgTetragonReady) Cast(o interface{}) notify.Message {
	return &MsgTetragonReady{}
}
