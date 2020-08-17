// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package cmd

import (
	"io"
	"os"

	"github.com/ethersphere/resolver/pkg/server"
	"github.com/ethersphere/resolver/pkg/version"
	"github.com/spf13/cobra"
)

type (
	Command = command
	Option  = option
)

var (
	NewCommand = newCommand
)

// RootCmd will return the internal root Cobra command.
func (cmd *Command) RootCmd() *cobra.Command {
	return cmd.root
}

// IntChan will return the command interrupt channel.
func (cmd *Command) IntChan() chan os.Signal {
	return cmd.intChan
}

// WithArgs will set the args for the command and pass it to the root Cobra
// command.
func WithArgs(args ...string) func(*Command) {
	return func(cmd *Command) {
		cmd.root.SetArgs(args)
	}
}

// WithBaseConfigDir will set the optional base configuration directory for the
// command.
func WithBaseConfigDir(baseConfigDir string) func(*Command) {
	return func(cmd *Command) {
		cmd.baseConfigDir = baseConfigDir
	}
}

// // WithConfigPath will override the path to the config file.
// func WithConfigPath(configPath string) func(*Command) {
// 	return func(cmd *Command) {
// 		cmd.configPath = configPath
// 	}
// }

// // WithCmdErr will override the standard error output for the command.
// func WithCmdErr(w io.Writer) func(*Command) {
// 	return func(cmd *Command) {
// 		cmd.root.SetErr(w)
// 	}
// }

// WithCmdOut will override the standard output for the command.
func WithCmdOut(w io.Writer) func(*Command) {
	return func(cmd *Command) {
		cmd.root.SetOut(w)
	}
}

// WithServerService will override the Server service implementation.
func WithServerService(s server.Service) func(*Command) {
	return func(cmd *Command) {
		cmd.services.server = s
	}
}

// WithVersionService will override the Version service implementation.
func WithVersionService(v version.Service) func(*Command) {
	return func(cmd *Command) {
		cmd.services.version = v
	}
}

// GetServerService will return the current Server service.
func (cmd *Command) GetServerService() server.Service {
	return cmd.services.server
}
