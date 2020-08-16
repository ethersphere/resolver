// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

// Client is a Resolver interface that can connect/disconnect to an external
// name resolution service via an edpoint.
type Client interface {
	Interface
	Connect(endpoint string) error
}
