// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/paxthemax/resolver/cmd/resolver-cli/cmd"
)

func TestMain(m *testing.M) {
	d, err := ioutil.TempDir("", "resolver-cmd-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	code := m.Run()

	if err := os.RemoveAll(d); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(code)
}

func newCommand(t *testing.T, opts ...cmd.Option) (c *cmd.Command) {
	t.Helper()

	c, err := cmd.NewCommand(opts...)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
