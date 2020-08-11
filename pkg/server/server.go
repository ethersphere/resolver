// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"net"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/ethersphere/resolver/pkg/api"
)

// Make sure Server implements Service.
var _ Service = (*Server)(nil)

// Service is the interface for the server package.
type Service interface {
	Address() string
	Serve() error
}

// Server wraps a HTTP server implementation, implementing the server package
// interface.
type Server struct {
	impl   http.Server
	logger *logrus.Logger
}

// Options are the Server options.
type Options struct {
	Addr   string
	Logger *logrus.Logger
}

// New creates and instantiates a new Server
func New(opts Options) *Server {
	srv := &Server{
		impl: http.Server{
			Addr: opts.Addr,
		},
		logger: opts.Logger,
	}

	return srv
}

// Address is the configured address the HTTP server will listen on.
func (s *Server) Address() string {
	return s.impl.Addr
}

// Serve will start the HTTP server implementation.
func (s *Server) Serve() error {
	s.logger.Infof("starting server on address %q", s.impl.Addr)

	apiService := api.New()
	s.impl.Handler = apiService

	// Create a net listener for the provided address.
	apiListener, err := net.Listen("tcp", s.impl.Addr)
	if err != nil {
		return err
	}

	// Serve the API.
	return s.impl.Serve(apiListener)
}
