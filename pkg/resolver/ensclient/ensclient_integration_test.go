// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build integraton

package ensclient_test

import (
	"strings"
	"testing"

	"github.com/ethersphere/resolver/pkg/resolver/ensclient"
)

func TestEnsclientIntegration(t *testing.T) {
	defaultEndpoint := "https://cloudflare-eth.com"

	testCases := []struct {
		desc            string
		endpoint        string
		name            string
		wantAdr         string
		wantFailConnect bool
		wantFailResolve bool
	}{
		{
			desc:            "bad ethclient endpoint",
			endpoint:        "fail",
			wantFailConnect: true,
		},
		{
			desc:            "bad name",
			name:            "iamaverybadname",
			wantFailResolve: true,
		},
		{
			desc:    "regular name",
			name:    "nickjohnson.eth",
			wantAdr: "0xb8c2c29ee19d8307cb7255e1cd9cbde883a267d5",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.endpoint == "" {
				tC.endpoint = defaultEndpoint
			}

			eC := ensclient.NewClient()

			err := eC.Connect(tC.endpoint)
			if err != nil {
				if !tC.wantFailConnect {
					t.Fatalf("failed to connect: %v", err)
				}
				return
			}

			adr, err := eC.Resolve(tC.name)
			if err != nil {
				if !tC.wantFailResolve {
					t.Fatalf("failed to resolve name: %v", err)
				}
				return
			}

			want := strings.ToLower(tC.wantAdr)
			got := strings.ToLower(adr.String())
			if got != want {
				t.Errorf("bad addr: got %q, want %q", got, want)
			}
		})
	}
}
