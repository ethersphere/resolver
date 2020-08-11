// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version

import "github.com/ethersphere/resolver"

// Make sure Version implements Service.
var _ Service = (*Version)(nil)

// Service is the interface for the version package.
type Service interface {
	String() string
}

// Version is the implementation of the Version service.
type Version struct{}

// New returns a new Version service.
func New() *Version {
	return &Version{}
}

// Version returns the current semantic version information.
func (v *Version) String() string {
	return resolver.Version
}
