// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package main

import (
	"github.com/go-faster/tetragon/cmd/tetra/loglevel"
	"github.com/go-faster/tetragon/cmd/tetra/tracingpolicy"
	"github.com/spf13/cobra"
)

func addCommands(rootCmd *cobra.Command) {
	addBaseCommands(rootCmd)
	rootCmd.AddCommand(tracingpolicy.New())
	rootCmd.AddCommand(loglevel.New())
}
