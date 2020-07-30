# Swarm name resolver

This repo contains the Swarm name resolver library/client. The name resolver can be used to attempt to resolve a name (URI) into an hash, and vice-versa.

This project is a work in progress.

## Usage:

`resolver-cli resolve <name>` - resolves an ENS name to an address string

NOTE: name MUST be a valid ENS name (must end with `.eth`).

## Structure:

- cmd/resolver-cli - a simple resolver command line client
- pkg/resolver - Multi resolver library implementation
- pkg/ens - ENS resolver abstraction

## TODO:

- Determine CLI format
- Decide what level of validation is required
- Add resolver endpoints to configuration (currently using cloudflare-eth)

