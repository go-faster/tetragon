// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package eventcache

import (
	"errors"
	"time"

	"github.com/go-faster/tetragon/api/v1/tetragon"
	"github.com/go-faster/tetragon/pkg/ktime"
	"github.com/go-faster/tetragon/pkg/metrics/errormetrics"
	"github.com/go-faster/tetragon/pkg/metrics/eventcachemetrics"
	"github.com/go-faster/tetragon/pkg/metrics/mapmetrics"
	"github.com/go-faster/tetragon/pkg/option"
	"github.com/go-faster/tetragon/pkg/process"
	"github.com/go-faster/tetragon/pkg/reader/node"
	"github.com/go-faster/tetragon/pkg/reader/notify"
	"github.com/go-faster/tetragon/pkg/server"
)

const (
	// garbage collection retries
	CacheStrikes = 15
	// garbage collection run interval
	EventRetryTimer = time.Second * 2
)

var (
	cache    *Cache
	nodeName string
)

type CacheObj struct {
	internal  *process.ProcessInternal
	event     notify.Event
	timestamp uint64
	startTime uint64
	color     int
	msg       notify.Message
}

type Cache struct {
	objsChan chan CacheObj
	done     chan bool
	cache    []CacheObj
	server   *server.Server
	dur      time.Duration
}

var (
	ErrFailedToGetPodInfo     = errors.New("failed to get pod info")
	ErrFailedToGetProcessInfo = errors.New("failed to get process info")
	ErrFailedToGetParentInfo  = errors.New("failed to get parent info")
)

// Generic internal lookup happens when events are received out of order and
// this event was handled before an exec event so it wasn't able to populate
// the process info yet.
func HandleGenericInternal(ev notify.Event, pid uint32, timestamp uint64) (*process.ProcessInternal, error) {
	internal, parent := process.GetParentProcessInternal(pid, timestamp)
	var err error

	if parent != nil {
		ev.SetParent(parent.GetProcessCopy())
	} else {
		errormetrics.ErrorTotalInc(errormetrics.EventCacheParentInfoFailed)
		err = ErrFailedToGetParentInfo
	}

	if internal != nil {
		ev.SetProcess(internal.GetProcessCopy())
	} else {
		errormetrics.ErrorTotalInc(errormetrics.EventCacheProcessInfoFailed)
		err = ErrFailedToGetProcessInfo
	}

	if err == nil {
		return internal, err
	}
	return nil, err
}

// Generic Event handler without any extra msg specific details or debugging
// so we only need to wait for the internal link to the process context to
// resolve PodInfo. This happens when the msg populates the internal state
// but that event is not fully populated yet.
func HandleGenericEvent(internal *process.ProcessInternal, ev notify.Event) error {
	p := internal.UnsafeGetProcess()
	if option.Config.EnableK8s && p.Pod == nil {
		errormetrics.ErrorTotalInc(errormetrics.EventCachePodInfoRetryFailed)
		return ErrFailedToGetPodInfo
	}

	ev.SetProcess(internal.GetProcessCopy())
	return nil
}

func (ec *Cache) handleEvents() {
	tmp := ec.cache[:0]
	for _, event := range ec.cache {
		var err error

		// If the process wasn't found before the Add(), likely because
		// the execve event was processed after this event, lets look it up
		// now because it should be available. Otherwise we have a valid
		// process and lets copy it across.
		if event.internal == nil {
			event.internal, err = event.msg.RetryInternal(event.event, event.startTime)
		}
		if err == nil {
			err = event.msg.Retry(event.internal, event.event)
		}
		if err != nil {
			event.color++
			if event.color < CacheStrikes {
				tmp = append(tmp, event)
				continue
			}
			if errors.Is(err, ErrFailedToGetProcessInfo) {
				eventcachemetrics.ProcessInfoError(notify.EventTypeString(event.event)).Inc()
			} else if errors.Is(err, ErrFailedToGetPodInfo) {
				eventcachemetrics.PodInfoError(notify.EventTypeString(event.event)).Inc()
			}
		}

		if event.msg.Notify() {
			processedEvent := &tetragon.GetEventsResponse{
				Event:    event.event.Encapsulate(),
				NodeName: nodeName,
				Time:     ktime.ToProto(event.timestamp),
			}

			ec.server.NotifyListeners(event.msg, processedEvent)
		}
	}
	ec.cache = tmp
}

func (ec *Cache) loop() {
	ticker := time.NewTicker(ec.dur)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			/* Every 'EventRetryTimer' walk the slice of events pending pod info. If
			 * an event hasn't completed its podInfo after two iterations send the
			 * event anyways.
			 */
			ec.handleEvents()
			mapmetrics.MapSizeSet("eventcache", 0, float64(len(ec.cache)))

		case event := <-ec.objsChan:
			eventcachemetrics.EventCacheCount.Inc()
			ec.cache = append(ec.cache, event)

		case <-ec.done:
			return
		}
	}
}

// We handle two race conditions here one where the event races with
// a Tetragon execve event and the other -- much more common -- where we
// race with K8s watcher
// case 1 (execve race):
//
//	Its possible to receive this Tetragon event before the process event cache
//	has been populated with a Tetragon execve event. In this case we need to
//	cache the event until the process cache is populated.
//
// case 2 (k8s watcher race):
//
//	Its possible to receive an event before the k8s watcher receives the
//	podInfo event and populates the local cache. If we expect podInfo,
//	indicated by having a nonZero dockerID we cache the event until the
//	podInfo arrives.
func (ec *Cache) Needed(proc *tetragon.Process) bool {
	if proc == nil {
		return true
	}
	if option.Config.EnableK8s {
		if proc.Docker != "" && proc.Pod == nil {
			return true
		}
	}
	if proc.Binary == "" {
		return true
	}
	return false
}

func (ec *Cache) Add(internal *process.ProcessInternal,
	e notify.Event,
	t uint64,
	s uint64,
	msg notify.Message) {
	ec.objsChan <- CacheObj{internal: internal, event: e, timestamp: t, startTime: s, msg: msg}
}

func NewWithTimer(s *server.Server, dur time.Duration) *Cache {
	if cache != nil {
		cache.done <- true
	}

	cache = &Cache{
		objsChan: make(chan CacheObj),
		done:     make(chan bool),
		cache:    make([]CacheObj, 0),
		server:   s,
		dur:      dur,
	}
	nodeName = node.GetNodeNameForExport()
	go cache.loop()
	return cache
}

func New(s *server.Server) *Cache {
	return NewWithTimer(s, EventRetryTimer)
}

func Get() *Cache {
	return cache
}
