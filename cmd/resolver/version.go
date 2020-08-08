// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"github.com/ethersphere/resolver/pkg/version"
	"github.com/spf13/cobra"
)

func (c *command) initVersionCmd() {

	c.root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run:   c.versionRun,
	})
}

func (c *command) versionRun(cmd *cobra.Command, args []string) {
	// If no service is injected, initialize the Version service.
	if c.services.version == nil {
		c.services.version = version.New()
	}

	cmd.Println(c.services.version.String())
}
