// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ethersphere/resolver/cmd/resolver/cmd"
)

var (
	baseConfigDir string
	testTempDir   string
)

func TestMain(m *testing.M) {
	tempDir, err := ioutil.TempDir("", "resolver-cmd-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up test: %v", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	baseConfigDir = filepath.Join(tempDir, "config")
	if err := os.Mkdir(baseConfigDir, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up test: %v", err)
		os.Exit(1)
	}

	testTempDir = tempDir
	os.Exit(m.Run())
}

// newCommand will create a new test command and set all default test options
// required for testing such as the base config directory.
// Commands can be specified with overrides for default values in tests.
func newCommand(t *testing.T, opts ...cmd.Option) *cmd.Command {
	t.Helper()

	// Set the base configuration dir to the temp dir.
	basedirOpt := []cmd.Option{cmd.WithBaseConfigDir(baseConfigDir)}

	// Execute all generic
	c, err := cmd.NewCommand(append(basedirOpt, opts...)...)
	if err != nil {
		t.Fatal(err)
	}

	return c
}

func TestBaseTestDirs(t *testing.T) {
	_ = newCommand(t)

	_, err := os.Stat(baseConfigDir)
	if err != nil && os.IsNotExist(err) {
		t.Fatal(err)
	}
}

func TestDefaultConfigPath(t *testing.T) {
	c := newCommand(t,
		// Use a noop root command runner to ensure the configuration directory
		// is created.
		cmd.WitCmdNoopRun(),
		// We need to override args because tests manipulate STDIN!
		cmd.WithArgs(""),
		cmd.WithCmdOut(ioutil.Discard),
	)
	if err := c.Execute(); err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(baseConfigDir, "swarm", "resolver")
	got := c.ConfigPath()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
