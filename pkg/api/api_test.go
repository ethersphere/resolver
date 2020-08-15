// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ethersphere/bee/pkg/logging"
	"github.com/ethersphere/resolver/pkg/api"
	"resenje.org/web"
)

// newTestServer will create a server for testing the API.
func newTestServer(t *testing.T, logger logging.Logger) *http.Client {
	if logger == nil {
		logger = logging.New(ioutil.Discard, 0)
	}

	a := api.New(logger)
	ts := httptest.NewServer(a)

	t.Cleanup(ts.Close)

	return &http.Client{
		Transport: web.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			u, err := url.Parse(ts.URL + r.URL.String())
			if err != nil {
				return nil, err
			}
			r.URL = u
			return ts.Client().Transport.RoundTrip(r)
		}),
	}
}

func TestRootHandler(t *testing.T) {
	cl := newTestServer(t, nil)

	t.Run("/ GET", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := cl.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		got := string(b)
		want := api.RootHandlerMessage + "\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestRobotsTxtHandler(t *testing.T) {
	cl := newTestServer(t, nil)

	t.Run("/robots.txt GET", func(t *testing.T) {
		resp, err := cl.Get("/robots.txt")
		if err != nil {
			t.Fatal(err)
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		got := string(b)
		want := "User-agent: *\nDisallow: /"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
