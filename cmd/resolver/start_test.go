// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver_test

import (
	"testing"

	"github.com/ethersphere/resolver/cmd/resolver"
	"github.com/ethersphere/resolver/pkg/server/mock"
)

func TestStartCommand(t *testing.T) {
	testDefaultAddress := ":8080"

	testCases := []struct {
		desc     string
		args     []string
		wantAddr string
	}{
		{
			desc: "OK - no address",
			args: []string{"start"},
		},
		{
			desc:     "OK - address set",
			args:     []string{"start", "--address", ":6000"},
			wantAddr: ":6000",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.wantAddr == "" {
				tC.wantAddr = testDefaultAddress
			}
			svc := mock.New(tC.wantAddr)

			cmd := newCommand(t,
				resolver.WithArgs(tC.args...),
				resolver.WithServerService(svc),
			)

			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}

			want := svc
			got := cmd.GetServerService()

			// Test if server address is correctly set.
			wAddr := want.Address()
			gAddr := got.Address()
			if wAddr != gAddr {
				t.Errorf("server address mismatch: want %q, got %q", wAddr, gAddr)
			}
		})
	}
}
