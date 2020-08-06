// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/adrg/xdg"
	"github.com/ethersphere/bee/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

const (
	configFileName = "resolver.conf"

	optionNameVerbosity  = "verbosity"
	optionShortVerbosity = "v"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "resolver",
	Short: "Swarm resolver",
	Long: `Swarm resolver performs name resolution for the swarm bee.
	
Configuration is stored in the "resolver.conf" file. Configuration lookup is
performed in the following order:
	
- File path explicitly passed via --config switch
- $XDG_CONFIG_HOME/resolver.conf
- Local directory resolver.conf
	`,
	SilenceErrors: true,
	SilenceUsage:  true,
	Run:           rootRun,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
}

func rootRun(cmd *cobra.Command, args []string) {

	// If a config file is found, read it in. Log any errors later.
	err := viper.ReadInConfig()

	// Set the logger from the config. The default log level is warn.
	var logger logging.Logger
	switch v := strings.ToLower(viper.GetString(optionNameVerbosity)); v {
	case "0", "silent":
		logger = logging.New(ioutil.Discard, 0)
	case "1", "error":
		logger = logging.New(cmd.OutOrStdout(), logrus.ErrorLevel)
	case "2", "warn":
	default:
		logger = logging.New(cmd.OutOrStdout(), logrus.WarnLevel)
	case "3", "info":
		logger = logging.New(cmd.OutOrStdout(), logrus.InfoLevel)
	case "4", "debug":
		logger = logging.New(cmd.OutOrStdout(), logrus.DebugLevel)
	case "5", "trace":
		logger = logging.New(cmd.OutOrStdout(), logrus.TraceLevel)
	}

	// If there were no errors reading the config file inform the user:
	if err == nil {
		logger.Infof("using config file %v", viper.ConfigFileUsed())
	} else {
		logger.Warningf("no config file found")
	}

}

// Execute adds all child commands to the root command and sets all the flags
// appropriately. This only needs to be invoked one in main().
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent flags will be global for the entire application:

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to the config file")

	// Local flags will only run when the action is called directly:

	rootCmd.Flags().StringP(optionNameVerbosity, optionShortVerbosity, "info", "log verbosity level 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=trace")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Use YAML for config:
		viper.SetConfigType("yaml")
		// Search config in $XDG_CONFIG_HOME with the name ".resolver"
		viper.AddConfigPath(xdg.ConfigHome)
		// Search in the working directory:
		viper.AddConfigPath(".")
		viper.SetConfigName(configFileName)
	}

	viper.AutomaticEnv() // read in environment variables that match
}
