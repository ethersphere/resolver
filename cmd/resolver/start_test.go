// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver_test

import (
	"fmt"
	"testing"

	"github.com/ethersphere/resolver/cmd/resolver"
	"github.com/ethersphere/resolver/pkg/server/mock"
	"github.com/sirupsen/logrus"
)

func TestStartCommand(t *testing.T) {
	testDefaultAddress := ":8080"

	testCases := []struct {
		desc          string
		args          []string
		wantAddr      string
		wantVerbosity bool
		verbosity     logrus.Level
	}{
		{
			desc: "OK - no flags",
			args: []string{"start"},
		},
		{
			desc:     "OK - address set",
			args:     []string{"start", "--address", ":6000"},
			wantAddr: ":6000",
		},
		{
			desc:          "OK - verbosity string set",
			args:          []string{"start", "--verbosity", "error"},
			wantVerbosity: true,
			verbosity:     logrus.ErrorLevel,
		},
		{
			desc:          "OK - verbosity number set",
			args:          []string{"start", "-v", "5"},
			wantVerbosity: true,
			verbosity:     logrus.ErrorLevel,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.wantAddr == "" {
				tC.wantAddr = testDefaultAddress
			}
			if !tC.wantVerbosity {
				tC.verbosity = logrus.InfoLevel
			}

			svc := mock.New(
				mock.WithAddress(tC.wantAddr),
				mock.WithLogLevel(tC.verbosity),
			)

			cmd := newCommand(t,
				resolver.WithArgs(tC.args...),
				resolver.WithServerService(svc),
			)

			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}

			want := svc
			got, ok := cmd.GetServerService().(*mock.Server)
			if !ok {
				t.Fatalf("test error: could not convert mock server")
			}

			// Test if server address is correctly set.
			wAddr := want.Address()
			gAddr := got.Address()
			if wAddr != gAddr {
				t.Errorf("server address mismatch: want %q, got %q", wAddr, gAddr)
			}

		})
	}
}

func TestStartCommandVerbosity(t *testing.T) {
	testCases := []struct {
		level     string
		wantLevel logrus.Level
	}{
		{
			level:     "0",
			wantLevel: 0,
		},
		{
			level:     "silent",
			wantLevel: 0,
		},
		{
			level:     "1",
			wantLevel: 1,
		},
		{
			level:     "error",
			wantLevel: 1,
		},
		{
			level:     "2",
			wantLevel: 2,
		},
		{
			level:     "warn",
			wantLevel: 2,
		},
		{
			level:     "3",
			wantLevel: 3,
		},
		{
			level:     "debug",
			wantLevel: 3,
		},
		{
			level:     "4",
			wantLevel: 4,
		},
		{
			level:     "trace",
			wantLevel: 4,
		},
	}
	for _, tC := range testCases {
		tDesc := fmt.Sprintf("set verbosity %s", tC.level)
		t.Run(tDesc, func(t *testing.T) {
			args := []string{"start", "--verbosity", tC.level}

			svc := mock.New(
				mock.WithLogLevel(tC.wantLevel),
			)

			cmd := newCommand(t,
				resolver.WithArgs(args...),
				resolver.WithServerService(svc),
			)

			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}

			want := svc
			got, ok := cmd.GetServerService().(*mock.Server)
			if !ok {
				t.Fatalf("test error: could not convert mock server")
			}

			if want.LogLevel != got.LogLevel {
				t.Errorf("server log verbosity mismatch: want %q, got %q", want.LogLevel, got.LogLevel)
			}
		})

	}
}
