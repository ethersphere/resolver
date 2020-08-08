// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ethersphere/resolver/cmd/resolver"
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
func newCommand(t *testing.T, opts ...resolver.Option) *resolver.Command {
	t.Helper()

	// Set the base configuration dir to the temp dir.
	basedirOpt := []resolver.Option{resolver.WithBaseConfigDir(baseConfigDir)}

	// Execute all generic
	cmd, err := resolver.NewCommand(append(basedirOpt, opts...)...)
	if err != nil {
		t.Fatal(err)
	}

	return cmd
}

func TestBaseTestDirs(t *testing.T) {
	_ = newCommand(t)

	_, err := os.Stat(baseConfigDir)
	if err != nil && os.IsNotExist(err) {
		t.Fatal(err)
	}
}
