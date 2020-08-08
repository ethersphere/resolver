// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import "net/http"

// Make sure Server implements Interface.
var _ Interface = (*Server)(nil)

// Interface is the interface for the server.
type Interface interface {
	Address() string
	Serve() error
}

// Server wraps a HTTP server implementation, implementing the server package
// interface.
type Server struct {
	impl http.Server
}

// New creates and instantiates a new Server
func New(addr string) *Server {
	srv := &Server{
		impl: http.Server{
			Addr: addr,
		},
	}

	return srv
}

// Address is the configured address the HTTP server will listen on.
func (s *Server) Address() string {
	return s.impl.Addr
}

// Serve will start the HTTP server implementation.
func (s *Server) Serve() error {
	return s.impl.ListenAndServe()
}
