// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"github.com/ethersphere/resolver/pkg/server"
	"github.com/spf13/cobra"
)

const (
	defaultAddress = ":8080"

	optionNameAddress = "address"
)

func (c *command) initStartCommand() {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the name resolution server",
		RunE:  c.startRunE,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.config.BindPFlags(cmd.Flags())
		},
	}

	cmd.Flags().String(optionNameAddress, defaultAddress, "Address for the server to listen on")
	c.root.AddCommand(cmd)
}

func (c *command) startRunE(cmd *cobra.Command, args []string) error {
	// If no service is injected, create a new Server service.
	if c.services.server == nil {
		c.services.server = server.New(c.config.GetString(optionNameAddress))
	}

	// Run the server loop.
	return c.services.server.Serve()
}
