// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mock

import (
	"fmt"

	"github.com/ethersphere/bee/pkg/swarm"
	"github.com/ethersphere/resolver/pkg/resolver"
)

// Address is the swarm bzz address.
type Address = swarm.Address

// Assure mock Resolver implements the Resolver interface.
var _ resolver.Interface = (*Resolver)(nil)

// Resolver is the mock Resolver implementation.
type Resolver struct {
	resolveFunc func(string) (Address, error)
}

// Option function sets the option on the mock Resolver.
type Option func(*Resolver)

// NewResolver will create a new mock Resolver.
func NewResolver(opts ...Option) resolver.Interface {
	r := &Resolver{}

	// Apply all options.
	for _, o := range opts {
		o(r)
	}

	return r
}

// WithResolveFunc will override the Resolve function implementation.
func WithResolveFunc(f func(string) (Address, error)) Option {
	return func(r *Resolver) {
		r.resolveFunc = f
	}
}

// Resolve implements the Resolver interface.
func (r *Resolver) Resolve(name string) (Address, error) {
	if r.resolveFunc != nil {
		return r.resolveFunc(name)
	}
	return Address{}, fmt.Errorf("not implemented")
}
