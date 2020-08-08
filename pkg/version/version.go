// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version

import "github.com/ethersphere/resolver"

// Make sure Version implements Interface.
var _ Interface = (*Version)(nil)

// Interface is the interface for the Version package.
type Interface interface {
	String() string
}

// Version returns the semantic version information of the package.
type Version struct {
}

// New returns a new Version service.
func New() *Version {
	return &Version{}
}

// Version returns the current semantic version information.
func (v *Version) String() string {
	return resolver.Version
}
