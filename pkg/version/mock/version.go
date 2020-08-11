// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mock

import "github.com/ethersphere/resolver/pkg/version"

// Make sure mock Version implements Version interface.
var _ version.Service = (*Version)(nil)

// Version is the mock versioning implementation.
type Version struct {
	val string
}

// New creates a new mock Version and sets the version string.
func New(val string) *Version {
	return &Version{
		val: val,
	}
}

// String implements the Version interface.
func (v *Version) String() string {
	return v.val
}
