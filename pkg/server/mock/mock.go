// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mock

import (
	"github.com/ethersphere/resolver/pkg/server"
	"github.com/sirupsen/logrus"
)

// Make sure mock Server implements server Interface.
var _ server.Interface = (*Server)(nil)

// Server is the mock server implementation.
type Server struct {
	addr     string
	LogLevel logrus.Level
	err      error
}

// Option applies an option to Server.
type Option func(*Server)

// New creates a new mock Server.
func New(opts ...Option) *Server {
	srv := &Server{}

	// Apply all options to the server:
	for _, o := range opts {
		o(srv)
	}

	return srv
}

// WithError sets the eror message returned by Serve.
func WithError(err error) Option {
	return func(s *Server) {
		s.err = err
	}
}

// WithAddress sets the configured mock server address.
func WithAddress(addr string) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

// WithLogLevel sets the configured mock server logging verbosity level.
func WithLogLevel(logLevel logrus.Level) Option {
	return func(s *Server) {
		s.LogLevel = logLevel
	}
}

// Address is the configured mock server addresss.
func (s *Server) Address() string {
	return s.addr
}

// Serve is the mock serve implementation
func (s *Server) Serve() error {
	return s.err
}
