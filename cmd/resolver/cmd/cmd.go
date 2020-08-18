// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ethersphere/resolver/pkg/server"
	"github.com/ethersphere/resolver/pkg/version"
)

func init() {
	cobra.EnableCommandSorting = false
}

// services contains all the services that defined commands can invoke.
type services struct {
	server  server.Service
	version version.Service
}

// command is the wrapper around a Cobra Command that contains everything that
// is required to unit test test said command, such as the reference to the
// root command,
type command struct {
	services      services
	baseConfigDir string
	config        *viper.Viper
	configDir     string
	root          *cobra.Command
	rootRunE      func(*cobra.Command, []string) error
	intChan       chan os.Signal
}

// option is used to apply optional parameters to a command.
type option func(*command)

// newCommand will create a new command with possible optional parameters
// applied.
func newCommand(opts ...option) (cmd *command, err error) {
	cmd = &command{
		root: &cobra.Command{
			Use:   "resolver",
			Short: "Swarm address resolver",
			Long: `Swarm resolver performs name resolution for the swarm bee.

Configuration is stored in the "config.yaml" file. Configuration lookup is
performed in the following order:
				
- File path explicitly passed via --config switch
- swarm/resolver/config.yaml in the standard user config directory 
    - eg. "~/.config/swarm/resolver/config.yaml" on Linux
- Working directory config.yaml`,
			SilenceErrors: true,
			SilenceUsage:  true,
			PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
				return cmd.initConfig()
			},
		},
	}

	// Apply all options to the command.
	for _, o := range opts {
		o(cmd)
	}

	cmd.initGlobalFlags()

	// Initialize all commands.
	cmd.initStartCommand()
	cmd.initVersionCmd()

	// Override for the root command RunE (used for tests).
	if cmd.rootRunE != nil {
		cmd.root.RunE = cmd.rootRunE
	}

	return cmd, nil
}

// Execute will execute the Cobra rootCommand.
func (cmd *command) Execute() error {
	return cmd.root.Execute()
}

// Execute adds all child commands to the root command and sets all the flags
// appropriately. This only needs to be invoked once in main().
func Execute() error {
	c, err := newCommand()
	if err != nil {
		return err
	}

	return c.Execute()
}

// initGlobalFlags will initialize all persistent flags on the root Cobra
// command.
func (cmd *command) initGlobalFlags() {
	flags := cmd.root.PersistentFlags()

	flags.StringVar(&cmd.configDir, "config", "", "path to the config file")
}

// initConfig will load the configuration from the config file and environment
// and load it into Viper.
func (cmd *command) initConfig() error {
	config := viper.New()
	configName := "config.yaml"

	// Set the system base config directory (eg. $XDG_CONFIG_HOME on linux).
	if cmd.baseConfigDir == "" {
		cmd.baseConfigDir = xdg.ConfigHome
	}

	// Set the config path:
	if cmd.configDir != "" {
		// Config file was explicitly passed via a flag, use it.
		config.SetConfigFile(cmd.configDir)
	} else {
		// Set the name and type of config file to search for.
		config.SetConfigName(configName)
		config.SetConfigType("yaml")

		// Obtain default config directory. If directory cannot be created, do
		// not include it in the config search path.
		// Path should default to "$XDG_CONFIG_HOME/swarm/resolver/config.yaml"
		configPath, err := createPath(cmd.baseConfigDir, "swarm", "resolver", "")
		if err == nil {
			config.AddConfigPath(configPath)
		}

		cmd.configDir = configPath
	}

	// Load the environment:
	config.SetEnvPrefix("resolver")
	config.AutomaticEnv() // Auto load all keys that match the prefix.
	config.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Load the config file:
	if err := config.ReadInConfig(); err != nil {
		// Do not return an error if config file is not found.
		var e viper.ConfigFileNotFoundError
		if !errors.As(err, &e) && !os.IsNotExist(err) {
			return err
		}
	}

	// Bind all flags to Viper.
	if err := config.BindPFlags(cmd.root.Flags()); err != nil {
		return err
	}

	cmd.config = config
	return nil
}

func createPath(name ...string) (string, error) {
	path := filepath.Join(name...)
	err := os.MkdirAll(path, os.ModeDir|0700)

	return path, err
}
