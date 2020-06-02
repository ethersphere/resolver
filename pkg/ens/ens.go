// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ens

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ens "github.com/wealdtech/go-ens/v3"
)

// Resolver is the ENS resolver client.
type Resolver struct {
	endpoint string
	client   *ethclient.Client
}

// NewResolver constructs an ENS resolver with a given endpoint.
func NewResolver(endpoint string) (res *Resolver) {
	res = &Resolver{
		endpoint: endpoint,
	}
	return res
}

// Connect will attempt to connect to the resolver endpoint, and initialize a client.
func (res *Resolver) Connect() (err error) {
	c, err := ethclient.Dial(res.endpoint)
	if err != nil {
		return err
	}

	res.client = c
	_, err = ens.PublicResolverAddress(c)
	return err
}

// Endpoint return the resolver endpoint.
func (res *Resolver) Endpoint() (endpoint string) {
	return res.endpoint
}

// IsConnected returns true if the resolver is connected to the endpoint.
func (res *Resolver) IsConnected() (ok bool) {
	return res.client != nil
}

// Resolve will try to resolve an ENS name into an address.
func (res *Resolver) Resolve(name string) (adr common.Address, err error) {
	adr, err = ens.Resolve(res.client, name)
	return
}
