// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ethersphere/bee/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ethersphere/resolver/pkg/server"
)

const (
	shutdownTimeout = 15 * time.Second

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

	// Initialize the signal channel.
	c.intChan = make(chan os.Signal, 1)
	signal.Notify(c.intChan, syscall.SIGINT, syscall.SIGTERM)

	cmd.Flags().String(optionNameAddress, defaultAddress, "address for the server to listen on")
	cmd.Flags().StringP(optionNameVerbosity, "v", defaultVerbosityLevel, "log verbosity level")
	c.root.AddCommand(cmd)
}

func (c *command) startRunE(cmd *cobra.Command, args []string) error {

	var logger logging.Logger
	switch v := strings.ToLower(c.config.GetString(optionNameVerbosity)); v {
	case "0", "silent":
		logger = logging.New(ioutil.Discard, 0)
	case "1", "error":
		logger = logging.New(cmd.OutOrStdout(), logrus.ErrorLevel)
	case "2", "warn":
		logger = logging.New(cmd.OutOrStdout(), logrus.WarnLevel)
	case "3", "info":
		logger = logging.New(cmd.OutOrStdout(), logrus.InfoLevel)
	case "4", "debug":
		logger = logging.New(cmd.OutOrStdout(), logrus.DebugLevel)
	case "5", "trace":
		logger = logging.New(cmd.OutOrStdout(), logrus.TraceLevel)
	default:
		return fmt.Errorf("Cannot set verbosity to %q", v)
	}

	// If no service is injected, create a new Server service.
	if c.services.server == nil {
		c.services.server = server.New(server.Options{
			Addr:   c.config.GetString(optionNameAddress),
			Logger: logger,
		})
	}

	// Run the resolver server loop.
	if err := c.services.server.Serve(); err != nil {
		return err
	}

	// Wait for interrupts from SIGINT/SIGTERM.
	// Block the main goroutine until interrupted.
	sig := <-c.intChan

	logger.Debugf("received signal: %v", sig)
	logger.Infof("server is shutting down")

	// Handle graceful shutdown.
	done := make(chan struct{})
	go func() {
		defer close(done)

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := c.services.server.Shutdown(ctx); err != nil {
			logger.Errorf("failed to shut down server: %v", err)
		}
	}()

	// Allow the forced server close by catching another signal.
	select {
	case sig := <-c.intChan:
		logger.Debugf("received signal: %v", sig)
		logger.Infof("interrupt received, terminating")
	case <-done:
	}

	return nil
}
