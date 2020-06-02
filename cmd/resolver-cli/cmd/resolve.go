package cmd

import (
	"errors"
	"strings"

	"github.com/paxthemax/resolver/pkg/resolver"
	"github.com/spf13/cobra"
)

// TODO: read this in from config.
const defaultENSEndpoint = "https://cloudflare-eth.com"

func (c *command) initResolveCmd() {
	c.root.AddCommand(&cobra.Command{
		Use:   "resolve",
		Short: "Resolve a name to a hash",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Parse arguments:
			if len(args) != 1 {
				return errors.New("Expected single argument: name to resolve")
			}
			name := args[0]

			// TODO: handle this in the resolver code
			if !strings.HasSuffix(name, "eth") {
				return errors.New("Not a valid ENS name")
			}

			mr := resolver.NewMultiResolver()
			mr.RegisterENSResolver(defaultENSEndpoint)
			r, err := mr.ConnectENS()
			if err != nil {
				cmd.PrintErr(err)
				return err
			}

			adr, err := r.Resolve(name)
			if err != nil {
				cmd.PrintErr(err)
				return err
			}

			cmd.Println(adr.String())

			return nil
		},
	})
}
