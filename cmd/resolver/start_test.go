// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver_test

import (
	"bytes"
	"strconv"
	"syscall"
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
		wantErr       bool
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
		{
			desc:    "fail - bad verbosity string",
			args:    []string{"start", "--version", "asparagus"},
			wantErr: true,
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

			var out bytes.Buffer
			cmd := newCommand(t,
				resolver.WithArgs(tC.args...),
				resolver.WithCmdOut(&out),
				resolver.WithServerService(svc),
			)

			cmd.IntChan() <- syscall.SIGTERM
			gotErr := cmd.Execute()
			if gotErr != nil && !tC.wantErr {
				t.Fatalf("unexpected error: %v", gotErr)
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

func TestStartVerbosity(t *testing.T) {
	testVerbosity := func(svc *mock.Server, lvl string) {
		var out bytes.Buffer

		// Test with full verbosity command.
		args := []string{"start", "--verbosity", lvl}

		cmd := newCommand(t,
			resolver.WithArgs(args...),
			resolver.WithCmdOut(&out),
			resolver.WithServerService(svc),
		)

		cmd.IntChan() <- syscall.SIGTERM
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}

		// Test with short command.
		args = []string{"start", "-v", lvl}

		cmd = newCommand(t,
			resolver.WithArgs(args...),
			resolver.WithCmdOut(&out),
			resolver.WithServerService(svc),
		)

		cmd.IntChan() <- syscall.SIGTERM
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	}

	testSet := map[string]int{
		"silent": 0,
		"error":  1,
		"warn":   2,
		"info":   3,
		"debug":  4,
		"trace":  5,
	}

	for lvlStr, lvlNum := range testSet {
		l := (logrus.Level)(lvlNum)
		svc := mock.New(
			mock.WithLogLevel(l),
		)

		testVerbosity(svc, lvlStr)
		testVerbosity(svc, strconv.Itoa(lvlNum))
	}
}

func TestServerShutdown(t *testing.T) {
	args := []string{"start", "--verbosity", "info"}

	svc := mock.New(
		mock.WithLogLevel(logrus.InfoLevel),
	)

	var out bytes.Buffer
	cmd := newCommand(t,
		resolver.WithArgs(args...),
		resolver.WithCmdOut(&out),
		resolver.WithServerService(svc),
	)

	cmd.IntChan() <- syscall.SIGTERM
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

}
