// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version_test

import (
	"testing"

	"github.com/ethersphere/resolver"
	"github.com/ethersphere/resolver/pkg/version"
)

func TestVersion(t *testing.T) {
	v := version.New()
	want := resolver.Version
	got := v.String()
	if want != got {
		t.Errorf("Version mismatch: want %q, got %q", want, got)
	}
}
