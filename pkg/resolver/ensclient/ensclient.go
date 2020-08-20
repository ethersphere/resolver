// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ensclient

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wealdtech/go-ens/v3"

	"github.com/ethersphere/resolver/pkg/resolver"
)

// Make sure Client implements the resolver.Client interface.
var _ resolver.Client = (*Client)(nil)

type dialType func(string) (*ethclient.Client, error)
type resolveType func(bind.ContractBackend, string) (common.Address, error)

// Client is a name resolution client that can connect to ENS/RNS via an
// Ethereum or RSK node endpoint.
type Client struct {
	Endpoint  string
	ethCl     *ethclient.Client
	dialFn    dialType
	resolveFn resolveType
}

// Option is a function that applies an option to a Client.
type Option func(*Client)

func wrapDial(ep string) (*ethclient.Client, error) {

	// Open a connection to the ethereum node through the endpoint.
	cl, err := ethclient.Dial(ep)
	if err != nil {
		return nil, err
	}

	// Ensure the ENS resolver contract is deployed on the network we are now
	// connected to.
	if _, err := ens.PublicResolverAddress(cl); err != nil {
		return nil, err
	}

	return cl, nil
}

// NewClient will return a new Client.
func NewClient(opts ...Option) *Client {
	c := &Client{
		dialFn:    wrapDial,
		resolveFn: ens.Resolve,
	}

	// Apply all options to the Client.
	for _, o := range opts {
		o(c)
	}

	return c
}

// Connect implements the resolver.Client interface.
func (c *Client) Connect(ep string) error {
	if c.dialFn == nil {
		return errors.New("no dial function implementation")
	}

	ethCl, err := c.dialFn(ep)
	if err != nil {
		return err
	}

	c.Endpoint = ep
	c.ethCl = ethCl
	return nil
}

// Resolve implements the resolver.Client interface.
func (c *Client) Resolve(name string) (resolver.Address, error) {
	if c.resolveFn == nil {
		return resolver.Address{}, errors.New("no resolve function implementation")
	}
	return c.resolveFn(c.ethCl, name)
}
