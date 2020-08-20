// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ensclient

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wealdtech/go-ens/v3"

	"github.com/ethersphere/bee/pkg/swarm"
	"github.com/ethersphere/resolver/pkg/resolver"
)

// Address is the swarm bzz address.
type Address = resolver.Address

// Make sure Client implements the resolver.Client interface.
var _ resolver.Client = (*Client)(nil)

type dialFn func(string) (*ethclient.Client, error)
type resolveFn func(bind.ContractBackend, string) (string, error)

// Client is a name resolution client that can connect to ENS/RNS via an
// Ethereum or RSK node endpoint.
type Client struct {
	Endpoint  string
	ethCl     *ethclient.Client
	dialFn    dialFn
	resolveFn resolveFn
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

func wrapResolve(backend bind.ContractBackend, name string) (string, error) {
	ethAdr, err := ens.Resolve(backend, name)
	return ethAdr.Hash().String(), err
}

// NewClient will return a new Client.
func NewClient(opts ...Option) *Client {
	c := &Client{
		dialFn:    wrapDial,
		resolveFn: wrapResolve,
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
func (c *Client) Resolve(name string) (Address, error) {
	if c.resolveFn == nil {
		return swarm.ZeroAddress, errors.New("no resolve function implementation")
	}
	hash, err := c.resolveFn(c.ethCl, name)
	if err != nil {
		return swarm.ZeroAddress, err
	}

	// Try and parse the raw address from the contract into a swarm address.
	return swarm.ParseHexAddress(hash)
}
