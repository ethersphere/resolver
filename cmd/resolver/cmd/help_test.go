// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ethersphere/resolver/cmd/resolver/cmd"
)

func TestRootCmdHelp(t *testing.T) {
	testCases := []struct {
		desc string
		args string
	}{
		{
			desc: "OK - no args",
			args: "",
		},
		{
			desc: "OK - passing '-h'",
			args: "-h",
		},
		{
			desc: "OK - passing '--help'",
			args: "--help",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var testOutBuf bytes.Buffer
			c := newCommand(t,
				cmd.WithArgs(tC.args),
				cmd.WithCmdOut(&testOutBuf),
			)

			if err := c.Execute(); err != nil {
				t.Fatal(err)
			}

			want := c.RootCmd().Long
			got := testOutBuf.String()

			if !strings.Contains(got, want) {
				t.Errorf("%q should contain %q", got, want)
			}
		})
	}
}
