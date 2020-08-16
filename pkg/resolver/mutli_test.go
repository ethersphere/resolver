// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethersphere/resolver/pkg/resolver"
	"github.com/ethersphere/resolver/pkg/resolver/mock"
)

type Address = common.Address

var (
	tErr  = errors.New("test resolve error")
	noAdr = Address{}
	tAdr1 = newAddr("deadbeef")
	tAdr2 = newAddr("beefdead")
)

func newAddr(s string) Address {
	return common.BytesToAddress([]byte(s))
}

func newResolver(adr Address, err error) resolver.Interface {
	return mock.NewResolver(
		mock.WithResolveFunc(
			func(string) (common.Address, error) { return adr, err }),
	)
}

func TestWithForceDefault(t *testing.T) {
	mr := resolver.NewMultiResolver(
		resolver.WithForceDefault(),
	)

	if !mr.ForceDefault {
		t.Error("did not set ForceDefault")
	}
}

func TestPushResolver(t *testing.T) {
	tld := ".tld"
	mr := resolver.NewMultiResolver()

	t.Run("error on bad tld", func(t *testing.T) {
		err := mr.PushResolver("invalid", &mock.Resolver{})
		want := resolver.ErrInvalidTLD
		if err != want {
			t.Errorf("bad error: got %v, want %v", err, want)
		}
	})

	t.Run("start empty", func(t *testing.T) {
		if mr.ChainCount(tld) > 0 {
			t.Fatal("not empty")
		}
	})

	t.Run("ok on tld", func(t *testing.T) {
		want := mock.NewResolver()
		err := mr.PushResolver(tld, want)
		if err != nil {
			t.Fatal(err)
		}
		got := mr.GetChain(tld)[0]
		if !reflect.DeepEqual(got, want) {
			t.Error("failed to push")
		}
	})
}

func TestPopResolver(t *testing.T) {
	tld := ".tld"
	mr := resolver.NewMultiResolver()

	t.Run("error on bad tld", func(t *testing.T) {
		err := mr.PopResolver("invalid")
		want := resolver.ErrInvalidTLD
		if err != want {
			t.Fatalf("bad error: got %v, want %v", err, want)
		}
	})

	t.Run("error on empty", func(t *testing.T) {
		err := mr.PopResolver(tld)
		want := resolver.ErrResolverChainEmpty
		if err != want {
			t.Fatalf("bad error: got %v, want %v", err, want)
		}
	})

	t.Run("ok on regular tld", func(t *testing.T) {
		mr.PushResolver(tld, mock.NewResolver())
		err := mr.PopResolver(tld)
		if err != nil {
			t.Error(err)
		}
		if mr.ChainCount(tld) > 0 {
			t.Error("failed to pop")
		}
	})
}

func TestResolve(t *testing.T) {
	testAdr := newAddr("aaaabbbbccccdddd")
	testAdrAlt := newAddr("ddddccccbbbbaaaa")

	newOKResolver := func(adr Address) resolver.Interface {
		return mock.NewResolver(
			mock.WithResolveFunc(func(_ string) (Address, error) {
				return adr, nil
			}),
		)
	}
	newErrResolver := func() resolver.Interface {
		return mock.NewResolver(
			mock.WithResolveFunc(func(name string) (Address, error) {
				err := fmt.Errorf("name resolution failed for %q", name)
				return Address{}, err
			}),
		)
	}

	testFixture := []struct {
		tld       string
		res       []resolver.Interface
		expectAdr Address
	}{
		{
			// Default chain:
			tld: "",
			res: []resolver.Interface{
				newOKResolver(testAdr),
			},
			expectAdr: testAdr,
		},
		{
			tld: ".tld",
			res: []resolver.Interface{
				newErrResolver(),
				newErrResolver(),
				newOKResolver(testAdr),
			},
			expectAdr: testAdr,
		},
		{
			tld: ".good",
			res: []resolver.Interface{
				newOKResolver(testAdr),
				newOKResolver(testAdrAlt),
			},
			expectAdr: testAdr,
		},
		{
			tld: ".empty",
		},
		{
			tld: ".fails",
			res: []resolver.Interface{
				newErrResolver(),
				newErrResolver(),
			},
		},
	}

	testCases := []struct {
		name    string
		wantAdr Address
		wantErr error
	}{
		{
			name:    "",
			wantAdr: testAdr,
		},
		{
			name:    "hello",
			wantAdr: testAdr,
		},
		{
			name:    "example.tld",
			wantAdr: testAdr,
		},
		{
			name:    ".tld",
			wantAdr: testAdr,
		},
		{
			name:    "get.good",
			wantAdr: testAdr,
		},
		{
			name:    "this.empty",
			wantErr: resolver.ErrResolverChainEmpty,
		},
		{
			name:    "this.fails",
			wantErr: fmt.Errorf("name resolution failed for %q", "this.fails"),
		},
	}

	// Load the test fixture.
	mr := resolver.NewMultiResolver()
	for _, tE := range testFixture {
		for _, r := range tE.res {
			if err := mr.PushResolver(tE.tld, r); err != nil {
				t.Fatal(err)
			}
		}
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			adr, err := mr.Resolve(tC.name)
			if err != nil {
				if tC.wantErr == nil {
					t.Fatalf("unexpected error: got %v", err)
				}
				if err.Error() != tC.wantErr.Error() {
					t.Fatalf("got %v, want %v", err, tC.wantErr)
				}
			}
			if adr != tC.wantAdr {
				t.Errorf("got %q, want %q", adr, tC.wantAdr)
			}
		})
	}
}
