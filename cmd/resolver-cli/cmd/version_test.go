// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"bytes"
	"testing"

	"github.com/ethersphere/resolver"
	"github.com/ethersphere/resolver/cmd/resolver-cli/cmd"
)

func TestVersionCmd(t *testing.T) {
	var outputBuffer bytes.Buffer
	if err := newCommand(t,
		cmd.WithArgs("version"),
		cmd.WithOutput(&outputBuffer),
	).Execute(); err != nil {
		t.Fatal(err)
	}

	want := resolver.Version + "\n"
	got := outputBuffer.String()
	if got != want {
		t.Errorf("Got version output %q, want %q", got, want)
	}
}
