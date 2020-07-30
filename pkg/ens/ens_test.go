// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ens_test

import (
	"encoding/hex"
	"testing"

	"github.com/ethersphere/resolver/pkg/ens"
)

const (
	cloudflareEndpoint = "https://cloudflare-eth.com"
)

func TestConnect(t *testing.T) {
	tests := map[string]struct {
		endpoint   string
		shouldFail bool
	}{
		"fails, empty": {
			endpoint:   "",
			shouldFail: true,
		},
		"fails, not an eth node": {
			endpoint:   "https://example.com",
			shouldFail: true,
		},
		"passes": {
			endpoint: cloudflareEndpoint,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res := ens.NewResolver(tc.endpoint)
			err := res.Connect()
			if tc.shouldFail {
				if err == nil {
					t.Fatalf("Test should fail")
				}
				return
			}

			if !res.IsConnected() {
				t.Fatalf("Function isConnected should return True")
			}
			wantEndpoint := res.Endpoint()
			gotEndpoint := tc.endpoint
			if wantEndpoint != gotEndpoint {
				t.Fatalf("Error in endpoint: want %v, got %v", wantEndpoint, gotEndpoint)
			}
		})
	}

}

func TestResolve(t *testing.T) {
	tests := map[string]struct {
		name       string
		want       string
		shouldFail bool
	}{
		"passes": {
			name: "nickjohnson.eth",
			want: "b8c2c29ee19d8307cb7255e1cd9cbde883a267d5",
		},
	}

	res := ens.NewResolver(cloudflareEndpoint)
	err := res.Connect()
	if err != nil {
		t.Fatalf("Error while connecting, err: %v", err)
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := res.Resolve(tc.name)
			if err != nil {
				t.Fatalf("Error in Resolve, err: %v", err)
			}
			got := hex.EncodeToString(result[:])
			if tc.want != got {
				t.Fatalf("Error in Resolve for name '%s': want %v, got %v", name, tc.want, got)
			}
		})
	}
}
