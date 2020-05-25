package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type command struct {
	root       *cobra.Command
	config     *viper.Viper
	configFile string
}

type option func(*command)

func newCommand(opts ...option) (c *command, err error) {
	c = &command{
		root: &cobra.Command{
			Use:           "resolver-cli",
			Short:         "Swarm name resolution client",
			SilenceErrors: true,
			SilenceUsage:  true,
		},
	}

	for _, o := range opts {
		o(c)
	}

	c.initGlobalFlags()

	if err := c.initConfig(); err != nil {
		return nil, err
	}

	c.initResolveCmd()
	c.initVersionCmd()

	return c, nil
}

func (c *command) Execute() (err error) {
	return c.root.Execute()
}

// Execute parses command line arguments and runs appropriate functions.
func Execute() (err error) {
	c, err := newCommand()
	if err != nil {
		return err
	}
	return c.Execute()
}

func (c *command) initGlobalFlags() {
	globalFlags := c.root.PersistentFlags()

	// Init config flag:
	globalFlags.StringVar(&c.configFile, "config", "", "Path to the resolver config file, default is ./config.yaml")
}

func (c *command) initConfig() (err error) {
	config := viper.New()

	config.SetConfigType("yaml")

	// Search for a config file:
	// If the --config flag is passed, use that file as config.
	// If no file is present, search the working directory.
	if c.configFile != "" {
		config.SetConfigFile(c.configFile)
	} else {
		config.AddConfigPath(".")
		config.SetConfigName("config")
	}

	config.SetEnvPrefix("RESOLVER")
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	err = config.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	c.config = config

	return nil
}
