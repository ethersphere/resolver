package resolver_test

import (
	"reflect"
	"testing"

	"github.com/paxthemax/resolver/pkg/ens"
	"github.com/paxthemax/resolver/pkg/resolver"
)

func makeENSResolvers(endpoints []string) (resolvers []resolver.Resolver) {
	for _, endpoint := range endpoints {
		resolvers = append(resolvers, ens.NewResolver(endpoint))
	}
	return resolvers
}

func TestRegisterResolvers(t *testing.T) {
	type inputs map[uint8]struct {
		endpoints []string
	}

	tests := map[string]inputs{
		"basic": {
			resolver.ResolverENS: {
				endpoints: []string{"test1", "test2", "test3"},
			},
		},
	}

	mr := resolver.NewMultiResolver()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ensInputs := tc[resolver.ResolverENS]
			for _, endpoint := range ensInputs.endpoints {
				mr.RegisterENSResolver(endpoint)
			}

			// TODO: add other resolver types.

			want := makeENSResolvers(ensInputs.endpoints)
			got := mr.GetENSResolvers()
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("Error creating ENS resolver: want %v, got %v", want, got)
			}

			// TODO: check other resolver types.
		})
	}
}

func TestResolverConnectENS(t *testing.T) {
	type endpoint struct {
		url   string
		fails bool
	}

	tests := map[string]struct {
		endpoints  []endpoint
		expectFail bool
	}{
		"fails, empty": {
			endpoints:  []endpoint{},
			expectFail: true,
		},
		"success, primary": {
			endpoints: []endpoint{
				{
					url: "https://cloudflare-eth.com",
				},
			},
		},
		"success, secondary": {
			endpoints: []endpoint{
				{
					url:   "fail1",
					fails: true,
				},
				{
					url:   "fail2",
					fails: true,
				},
				{
					url:   "https://cloudflare-eth.com",
					fails: false,
				},
			},
		},
		"success, primary only": {
			endpoints: []endpoint{
				{
					url:   "https://cloudflare-eth.com",
					fails: false,
				},
				{
					url:   "https://cloudflare-eth.com",
					fails: false,
				},
			},
		},
	}

	mr := resolver.NewMultiResolver()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for _, endpoint := range tc.endpoints {
				mr.RegisterENSResolver(endpoint.url)
			}
			// TODO; register other resolver types

			res, err := mr.ConnectENS()
			if tc.expectFail {
				if err == nil {
					t.Fatalf("Error: expected failure")
				}
				return
			}

			if err != nil {
				t.Fatalf("Error connecting to ENS, error: %v", err)
			}
			if !res.IsConnected() {
				t.Fatalf("Error connecting to ENS endpoint: IsConnected() is false")
			}

			// Find first endpoint that does not fail:
			wantEndpoint := ""
			for _, endpoint := range tc.endpoints {
				if !endpoint.fails {
					wantEndpoint = endpoint.url
				}
			}
			gotEndpoint := res.Endpoint()
			if gotEndpoint != wantEndpoint {
				t.Fatalf("Error connecting to ENS endpoint, want: %s, got: %s", wantEndpoint, gotEndpoint)
			}

		})
	}
}

func TestResolverConnectRNS(t *testing.T) {
	mr := resolver.NewMultiResolver()
	_, err := mr.ConnectRNS()
	if err == nil {
		t.Fatalf("Error: function not implemented")
	}
}
