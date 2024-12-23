// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

// Code generated by protoc-gen-go-tetragon. DO NOT EDIT

package helpers

import (
	fmt "fmt"
	tetragon "github.com/go-faster/tetragon/api/v1/tetragon"
	proto "google.golang.org/protobuf/proto"
)

// ResponseTypeString returns an event's type as a string
func ResponseTypeString(response *tetragon.GetEventsResponse) (string, error) {
	if response == nil {
		return "", fmt.Errorf("Response is nil")
	}

	event := response.Event
	if event == nil {
		return "", fmt.Errorf("Event is nil")
	}

	switch event.(type) {
	case *tetragon.GetEventsResponse_ProcessExec:
		return tetragon.EventType_PROCESS_EXEC.String(), nil
	case *tetragon.GetEventsResponse_ProcessExit:
		return tetragon.EventType_PROCESS_EXIT.String(), nil
	case *tetragon.GetEventsResponse_ProcessKprobe:
		return tetragon.EventType_PROCESS_KPROBE.String(), nil
	case *tetragon.GetEventsResponse_ProcessTracepoint:
		return tetragon.EventType_PROCESS_TRACEPOINT.String(), nil
	case *tetragon.GetEventsResponse_ProcessLoader:
		return tetragon.EventType_PROCESS_LOADER.String(), nil
	case *tetragon.GetEventsResponse_ProcessUprobe:
		return tetragon.EventType_PROCESS_UPROBE.String(), nil
	case *tetragon.GetEventsResponse_ProcessThrottle:
		return tetragon.EventType_PROCESS_THROTTLE.String(), nil
	case *tetragon.GetEventsResponse_ProcessLsm:
		return tetragon.EventType_PROCESS_LSM.String(), nil
	case *tetragon.GetEventsResponse_Test:
		return tetragon.EventType_TEST.String(), nil
	case *tetragon.GetEventsResponse_RateLimitInfo:
		return tetragon.EventType_RATE_LIMIT_INFO.String(), nil

	}
	return "", fmt.Errorf("Unhandled response type %T", event)
}

// ResponseGetProcess returns a GetEventsResponse's process if it exists
func ResponseGetProcess(response *tetragon.GetEventsResponse) *tetragon.Process {
	if response == nil {
		return nil
	}

	event := response.Event
	if event == nil {
		return nil
	}

	return ResponseInnerGetProcess(event)
}

// ResponseInnerGetProcess returns a GetEventsResponse inner event's process if it exists
func ResponseInnerGetProcess(event tetragon.IsGetEventsResponse_Event) *tetragon.Process {
	switch ev := event.(type) {
	case *tetragon.GetEventsResponse_ProcessExec:
		return ev.ProcessExec.Process
	case *tetragon.GetEventsResponse_ProcessExit:
		return ev.ProcessExit.Process
	case *tetragon.GetEventsResponse_ProcessKprobe:
		return ev.ProcessKprobe.Process
	case *tetragon.GetEventsResponse_ProcessTracepoint:
		return ev.ProcessTracepoint.Process
	case *tetragon.GetEventsResponse_ProcessUprobe:
		return ev.ProcessUprobe.Process
	case *tetragon.GetEventsResponse_ProcessLsm:
		return ev.ProcessLsm.Process
	case *tetragon.GetEventsResponse_ProcessLoader:
		return ev.ProcessLoader.Process

	}
	return nil
}

// ResponseGetProcessKprobe returns a GetEventsResponse's process if it exists
func ResponseGetProcessKprobe(response *tetragon.GetEventsResponse) *tetragon.ProcessKprobe {
	if response == nil {
		return nil
	}

	return response.GetProcessKprobe()
}

// ResponseGetParent returns a GetEventsResponse's parent process if it exists
func ResponseGetParent(response *tetragon.GetEventsResponse) *tetragon.Process {
	if response == nil {
		return nil
	}

	event := response.Event
	if event == nil {
		return nil
	}

	return ResponseInnerGetParent(event)
}

// ResponseInnerGetParent returns a GetEventsResponse inner event's parent process if it exists
func ResponseInnerGetParent(event tetragon.IsGetEventsResponse_Event) *tetragon.Process {
	switch ev := event.(type) {
	case *tetragon.GetEventsResponse_ProcessExec:
		return ev.ProcessExec.Parent
	case *tetragon.GetEventsResponse_ProcessExit:
		return ev.ProcessExit.Parent
	case *tetragon.GetEventsResponse_ProcessKprobe:
		return ev.ProcessKprobe.Parent
	case *tetragon.GetEventsResponse_ProcessTracepoint:
		return ev.ProcessTracepoint.Parent
	case *tetragon.GetEventsResponse_ProcessUprobe:
		return ev.ProcessUprobe.Parent
	case *tetragon.GetEventsResponse_ProcessLsm:
		return ev.ProcessLsm.Parent

	}
	return nil
}

// ResponseTypeMap returns a map from event field names (e.g. "process_exec") to corresponding
// protobuf messages (e.g. &tetragon.ProcessExec{}).
func ResponseTypeMap() map[string]proto.Message {
	return map[string]proto.Message{
		"process_exec":       &tetragon.ProcessExec{},
		"process_exit":       &tetragon.ProcessExit{},
		"process_kprobe":     &tetragon.ProcessKprobe{},
		"process_tracepoint": &tetragon.ProcessTracepoint{},
		"process_loader":     &tetragon.ProcessLoader{},
		"process_uprobe":     &tetragon.ProcessUprobe{},
		"process_throttle":   &tetragon.ProcessThrottle{},
		"process_lsm":        &tetragon.ProcessLsm{},
		"test":               &tetragon.Test{},
		"rate_limit_info":    &tetragon.RateLimitInfo{},
	}
}

// ProcessEventMap returns a map from event field names (e.g. "process_exec") to corresponding
// protobuf messages in a given tetragon.GetEventsResponse (e.g. response.GetProcessExec()).
func ProcessEventMap(response *tetragon.GetEventsResponse) map[string]any {
	return map[string]any{
		"process_exec":       response.GetProcessExec(),
		"process_exit":       response.GetProcessExit(),
		"process_kprobe":     response.GetProcessKprobe(),
		"process_tracepoint": response.GetProcessTracepoint(),
		"process_loader":     response.GetProcessLoader(),
		"process_uprobe":     response.GetProcessUprobe(),
		"process_throttle":   response.GetProcessThrottle(),
		"process_lsm":        response.GetProcessLsm(),
		"test":               response.GetTest(),
		"rate_limit_info":    response.GetRateLimitInfo(),
	}
}
