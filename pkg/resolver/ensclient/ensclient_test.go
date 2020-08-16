// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ensclient_test

import (
	"errors"
	"testing"

	"github.com/ethersphere/resolver/pkg/resolver"
	ec "github.com/ethersphere/resolver/pkg/resolver/ensclient"
)

func TestNewClient(t *testing.T) {
	cl := ec.NewClient()
	if cl.Endpoint != "" {
		t.Errorf("expected no endpoint set")
	}
}

func TestConnect(t *testing.T) {
	ep := "test"

	t.Run("no dial func error", func(t *testing.T) {
		c := ec.NewClient(
			ec.WithDialFunc(nil),
		)
		err := c.Connect(ep)
		if err == nil && err.Error() != "no dial function implementation" {
			t.Fatal("expected error")
		}
	})

	t.Run("connect error", func(t *testing.T) {
		c := ec.NewClient(
			ec.WithErrorDialFunc(errors.New("failed to connect")),
		)

		if err := c.Connect("test"); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("ok", func(t *testing.T) {
		c := ec.NewClient(
			ec.WithNoopDialFunc(),
		)

		if err := c.Connect(ep); err != nil {
			t.Fatal(err)
		}
		if c.Endpoint != ep {
			t.Errorf("bad endpoint: got %q, want %q", c.Endpoint, ep)
		}

	})
}

func TestResolve(t *testing.T) {
	name := "hello"

	t.Run("no resolve func error", func(t *testing.T) {
		c := ec.NewClient(
			ec.WithResolveFunc(nil),
		)
		_, err := c.Resolve("test")
		if err == nil && err.Error() != "no resolve function implementation" {
			t.Fatal("expected error")
		}
	})

	t.Run("resolve error", func(t *testing.T) {
		c := ec.NewClient(
			ec.WithNoopDialFunc(),
			ec.WithErrorResolveFunc(errors.New("resolve error")),
		)

		if err := c.Connect(name); err != nil {
			t.Fatal(err)
		}

		_, err := c.Resolve(name)
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("ok", func(t *testing.T) {
		c := ec.NewClient(
			ec.WithNoopDialFunc(),
			ec.WithAdrResolveFunc(resolver.Address{}),
		)

		if err := c.Connect(name); err != nil {
			t.Fatal(err)
		}

		adr, err := c.Resolve(name)
		if err != nil {
			t.Error(err)
		}
		want := (resolver.Address{}).String()
		got := adr.String()
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

}
