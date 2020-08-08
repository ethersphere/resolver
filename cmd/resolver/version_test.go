// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package resolver_test

import (
	"bytes"
	"testing"

	"github.com/ethersphere/resolver/cmd/resolver"
	"github.com/ethersphere/resolver/pkg/version/mock"
)

func TestVersionCmd(t *testing.T) {
	const testArgs = "version"
	const testVersionString = "TEST TEST TEST"
	var testOutBuf bytes.Buffer

	cmd := newCommand(t,
		resolver.WithArgs(testArgs),
		resolver.WithCmdOut(&testOutBuf),
		resolver.WithVersionService(mock.New(testVersionString)),
	)

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	want := testVersionString + "\n"
	got := testOutBuf.String()

	if want != got {
		t.Errorf("bad output: want %q, got %q", want, got)
	}
}
