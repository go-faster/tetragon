// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package grpc

import (
	"context"
	"sync"

	"github.com/cilium/hubble/pkg/cilium"
	"github.com/go-faster/tetragon/api/v1/tetragon"
	"github.com/go-faster/tetragon/pkg/eventcache"
	"github.com/go-faster/tetragon/pkg/logger"
	"github.com/go-faster/tetragon/pkg/metrics/eventmetrics"
	"github.com/go-faster/tetragon/pkg/option"
	"github.com/go-faster/tetragon/pkg/reader/node"
	"github.com/go-faster/tetragon/pkg/reader/notify"
	"github.com/go-faster/tetragon/pkg/sensors"
	"github.com/go-faster/tetragon/pkg/server"
	"github.com/sirupsen/logrus"
)

// ProcessManager maintains a cache of processes from tetragon exec events.
type ProcessManager struct {
	nodeName string
	Server   *server.Server
	// synchronize access to the listeners map.
	mux         sync.Mutex
	listeners   map[server.Listener]struct{}
	ciliumState *cilium.State
}

// NewProcessManager returns a pointer to an initialized ProcessManager struct.
func NewProcessManager(
	ctx context.Context,
	wg *sync.WaitGroup,
	ciliumState *cilium.State,
	manager *sensors.Manager,
) (*ProcessManager, error) {
	pm := &ProcessManager{
		nodeName:    node.GetNodeNameForExport(),
		ciliumState: ciliumState,
		listeners:   make(map[server.Listener]struct{}),
	}

	pm.Server = server.NewServer(ctx, wg, pm, manager)

	// Exec cache is always needed to ensure events have an associated Process{}
	eventcache.New(pm.Server)

	logger.GetLogger().WithField("enableCilium", option.Config.EnableCilium).WithFields(logrus.Fields{
		"enableK8s":         option.Config.EnableK8s,
		"enableProcessCred": option.Config.EnableProcessCred,
		"enableProcessNs":   option.Config.EnableProcessNs,
	}).Info("Starting process manager")
	return pm, nil
}

// Notify implements Listener.Notify.
func (pm *ProcessManager) Notify(event notify.Message) error {
	processedEvent := event.HandleMessage()
	if processedEvent != nil {
		pm.NotifyListener(event, processedEvent)
	}
	return nil
}

// Close implements Listener.Close.
func (pm *ProcessManager) Close() error {
	return nil
}

func (pm *ProcessManager) AddListener(listener server.Listener) {
	logger.GetLogger().WithField("getEventsListener", listener).Debug("Adding a getEventsListener")
	pm.mux.Lock()
	defer pm.mux.Unlock()
	pm.listeners[listener] = struct{}{}
}

func (pm *ProcessManager) RemoveListener(listener server.Listener) {
	logger.GetLogger().WithField("getEventsListener", listener).Debug("Removing a getEventsListener")
	pm.mux.Lock()
	defer pm.mux.Unlock()
	delete(pm.listeners, listener)
}

func (pm *ProcessManager) NotifyListener(original interface{}, processed *tetragon.GetEventsResponse) {
	pm.mux.Lock()
	defer pm.mux.Unlock()
	for l := range pm.listeners {
		l.Notify(processed)
	}
	eventmetrics.ProcessEvent(original, processed)
}
