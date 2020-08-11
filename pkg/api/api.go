// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import "net/http"

// Service is the interface for the api package.
type Service interface {
	http.Handler
}

// API is the implementation of the api Service.
// API wraps an http handler.
type API struct {
	http.Handler
}

// New will return a new API instance.
func New() Service {
	a := &API{}

	a.setupRouting()

	return a
}
