// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package observer

import (
	"os"

	"github.com/go-faster/tetragon/pkg/sensors"
)

func RemovePrograms(bpfDir, mapDir string) {
	sensors.UnloadAll(bpfDir)
	os.Remove(bpfDir)
	os.Remove(mapDir)
}
