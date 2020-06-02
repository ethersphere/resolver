// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"bytes"
	"testing"

	"github.com/paxthemax/resolver/cmd/resolver-cli/cmd"
)

func TestResolveCmd(t *testing.T) {
	tests := map[string]struct {
		args       []string
		want       string
		shouldFail bool
	}{
		"fails, no arguments": {
			shouldFail: true,
		},
		"fails, multiple arguments": {
			args:       []string{"no", "dice"},
			shouldFail: true,
		},
		"fails, name does not end with .eth": {
			args:       []string{"example.com"},
			shouldFail: true,
		},
		"passes, nickjohnson.eth": {
			args:       []string{"nickjohnson.eth"},
			want:       "0xb8c2C29ee19D8307cb7255e1Cd9CbDE883A267d5",
			shouldFail: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			args := append([]string{"resolve"}, tc.args...)
			var outputBuffer bytes.Buffer
			err := newCommand(t,
				cmd.WithArgs(args...),
				cmd.WithOutput(&outputBuffer),
			).Execute()
			if err != nil {
				if tc.shouldFail {
					return
				}
				t.Fatal(err)
			}

			want := tc.want + "\n"
			got := outputBuffer.String()
			if want != got {
				if tc.shouldFail {
					return
				}
				t.Fatalf("Eroor in resolve cmd: want %v, got %v", want, got)
			}
		})
	}
}
