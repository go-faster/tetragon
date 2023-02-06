// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package policyfilter

import (
	"fmt"
	"io"
	"runtime"

	"github.com/sirupsen/logrus"
)

// there is no way to have selective information level  per sub-system
// (see: https://github.com/cilium/cilium/issues/21002) so we define a flag and
// some helper functions here.

const (
	debugInfo = true
)

func initEmptylogger() logrus.FieldLogger {
	// NB: we could define a better empty logger, that also ignores WithField
	log := logrus.New()
	log.SetOutput(io.Discard)
	return log
}

var (
	emptyLogger = initEmptylogger()
)

func debugLog(log logrus.FieldLogger) logrus.FieldLogger {
	if !debugInfo {
		return emptyLogger
	}

	return log
}

func debugPodLogger(log logrus.FieldLogger, p *podInfo) logrus.FieldLogger {
	if !debugInfo {
		return emptyLogger
	}

	log = log.WithField("pod-id", p.id)
	cids := make([]string, 0, len(p.containers))
	for _, c := range p.containers {
		cids = append(cids, c.id)
	}
	log = log.WithField("container-ids", cids)
	return log
}

func (s *state) debugLogWithCallers(nCallers int) logrus.FieldLogger {
	if !debugInfo {
		return emptyLogger
	}

	log := s.log
	for i := 1; i <= nCallers; i++ {
		pc, _, _, ok := runtime.Caller(i)
		if !ok {
			return log
		}
		fn := runtime.FuncForPC(pc)
		key := fmt.Sprintf("caller-%d", i)
		log = log.WithField(key, fn.Name())
	}

	return log
}

func (s *state) Debug(args ...interface{}) {
	if debugInfo {
		s.log.Info(args...)
	} else {
		s.log.Debug(args...)
	}
}

func (s *state) Debugf(fmt string, args ...interface{}) {
	if debugInfo {
		s.log.Infof(fmt, args...)
	} else {
		s.log.Debugf(fmt, args...)
	}
}
