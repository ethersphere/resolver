// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethersphere/resolver/pkg/ens"
)

const (
	resolverENS uint8 = iota
	resolverRNS
)

// Resolver can resolve URIs to addresses, and reverse-resolve addresses to URIs
type Resolver interface {
	Connect() (err error)
	Endpoint() (endpoint string)
	IsConnected() (ok bool)
	Resolve(name string) (address common.Address, err error)
}

type resolverMap map[uint8][]Resolver

// MultiResolver contains all registered resolvers.
type MultiResolver struct {
	resolvers resolverMap
}

// NewMultiResolver creates a multi resolver.
func NewMultiResolver() (mr *MultiResolver) {
	return &MultiResolver{
		resolvers: make(resolverMap),
	}
}

func (mr *MultiResolver) register(resolverType uint8, endpoint string) {
	var resolver Resolver

	switch resolverType {
	case resolverENS:
		resolver = ens.NewResolver(endpoint)
	}
	mr.resolvers[resolverType] = append(mr.resolvers[resolverType], resolver)
}

// RegisterENSResolver registers a resolver to ENS with a given endpoint.
func (mr *MultiResolver) RegisterENSResolver(endpoint string) {
	mr.register(resolverENS, endpoint)
}

func (mr *MultiResolver) get(resolverType uint8) (resolvers []Resolver) {
	return mr.resolvers[resolverType]
}

// GetENSResolvers will return all registered ENS resolvers.
func (mr *MultiResolver) GetENSResolvers() (resolvers []Resolver) {
	return mr.get(resolverENS)
}

func (mr *MultiResolver) connect(resolverType uint8) (res Resolver, err error) {
	switch resolverType {
	case resolverENS:
		for _, r := range mr.GetENSResolvers() {
			if err := r.Connect(); err == nil {
				return r, nil
			}
		}
	case resolverRNS:
		return nil, errors.New("TODO: implement RNS support")
	}
	return nil, errors.New("unknown resolver requested")
}

// ConnectENS will attempt to connect to the ENS. MultiResolver will try all registered resolvers in the chain, until one connects.
// If no resolvers connect, the function will return an error.
func (mr *MultiResolver) ConnectENS() (res Resolver, err error) {
	return mr.connect(resolverENS)
}

// ConnectRNS will attempt to connect to the RNS. MultiResolver will try all registered resolvers in the chain, until one connects.
// If no resolvers connect, the function will return an error.
func (mr *MultiResolver) ConnectRNS() (res Resolver, err error) {
	return mr.connect(resolverRNS)
}
