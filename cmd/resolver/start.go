// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ethersphere/resolver/pkg/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	defaultAddress        = ":8080"
	defaultVerbosityLevel = "info"

	optionNameAddress   = "address"
	optionNameVerbosity = "verbosity"
)

func (c *command) initStartCommand() {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the name resolution server",
		Long: `Start the name resolution server.

The server will listen on the configured address.
Logging verbosity level can be provided as a number or a string:
	0: silent
	1: error
	2: warn
	3: info (default)
	4: trace`,
		RunE: c.startRunE,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.config.BindPFlags(cmd.Flags())
		},
	}

	cmd.Flags().String(optionNameAddress, defaultAddress, "address for the server to listen on")
	cmd.Flags().StringP(optionNameVerbosity, "v", defaultVerbosityLevel, "log verbosity level")
	c.root.AddCommand(cmd)
}

func (c *command) startRunE(cmd *cobra.Command, args []string) error {

	// TODO: configure logger verbosity, extract logging package, add metrics?
	logger := logrus.New()
	logger.SetOutput(cmd.OutOrStdout())

	switch v := strings.ToLower(c.config.GetString(optionNameVerbosity)); v {
	case "0", "silent":
		logger.SetLevel(0)
		logger.SetOutput(ioutil.Discard)
	case "1", "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "2", "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "3", "info":
		logger.SetLevel(logrus.InfoLevel)
	case "4", "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "5", "trace":
		logger.SetLevel(logrus.TraceLevel)
	default:
		return fmt.Errorf("Cannot set verbosity to %q", v)
	}

	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	// If no service is injected, create a new Server service.
	if c.services.server == nil {
		c.services.server = server.New(server.Options{
			Addr:   c.config.GetString(optionNameAddress),
			Logger: logger,
		})
	}

	// Run the resolver server loop.
	return c.services.server.Serve()
}
