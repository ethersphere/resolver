// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mock

import (
	"context"

	"github.com/ethersphere/resolver/pkg/server"
	"github.com/sirupsen/logrus"
)

// Make sure mock Server implements server Interface.
var _ server.Service = (*Server)(nil)

// Server is the mock server implementation.
type Server struct {
	addr        string
	LogLevel    logrus.Level
	err         error
	wasClosed   bool
	wasShutdown bool
	closeFn     func() error
	shutdownFn  func(context.Context) error
}

// Option applies an option to Server.
type Option func(*Server)

// New creates a new mock Server.
func New(opts ...Option) *Server {
	srv := &Server{}

	// Apply all options to the Server.
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

// WithCloseFn sets the mock Close function implementation.
func WithCloseFn(fn func() error) Option {
	return func(s *Server) {
		s.closeFn = fn
	}
}

// WithShutdownFn sets the mock Shutdown function implementation.
func WithShutdownFn(fn func(context.Context) error) Option {
	return func(s *Server) {
		s.shutdownFn = fn
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

// Close is the mock Close implementation.
func (s *Server) Close() error {
	s.wasClosed = true
	return nil
}

// Shutdown is the mock Shutdown implementation.
func (s *Server) Shutdown(ctx context.Context) error {
	s.wasShutdown = true
	return nil
}
