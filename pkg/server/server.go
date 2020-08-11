// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"context"
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
	Close() error
	Shutdown(ctx context.Context) error
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
	s.logger.Info("starting server")

	apiService := api.New()
	s.impl.Handler = apiService

	// Create a net listener for the provided address.
	l, err := net.Listen("tcp", s.impl.Addr)
	if err != nil {
		return err
	}

	// Start a goroutine to serve the API.
	go func(srv *Server, ln net.Listener) {
		srv.logger.Infof("server listening on %q", l.Addr().String())
		if err := srv.impl.Serve(ln); err != nil {
			srv.logger.Errorf("failed to serve on address %q: %v", l.Addr().String(), err)
		}
	}(s, l)

	return nil
}

// Close immediately closes any open connections and termintates the server.
// For a graceful server shutdown, use Shutdown.
func (s *Server) Close() error {
	return s.impl.Close()
}

// Shutdown will gracefully shut down the server without interrupting any
// active connections.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.impl.Shutdown(ctx)
}
