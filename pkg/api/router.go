// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"net/http"

	"github.com/ethersphere/bee/pkg/jsonhttp"
	"github.com/ethersphere/bee/pkg/logging"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"resenje.org/web"
)

const (
	robotsTxtMessage   = "User-agent: *\nDisallow: /"
	rootHandlerMessage = "Ethereum Swarm resolver"
)

// TODO: add versioning to the API.

// setupRouting will set all routes for the API.
func (a *API) setupRouting() {
	router := mux.NewRouter()

	// Handle 404.
	router.NotFoundHandler = http.HandlerFunc(jsonhttp.NotFoundHandler)

	// Handle the root path.
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, rootHandlerMessage)
	})

	// Disallow robots.
	router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, robotsTxtMessage)
	})

	// Add all middleware to the handler chain and set the top level handler.
	a.Handler = web.ChainHandlers(
		// Log all access to the API on the INFO level.
		logging.NewHTTPAccessLogHandler(a.logger, logrus.InfoLevel, "api access"),

		// TODO: Add compression handlers?
		// TODO: Add metrics handler?
		// TODO: Handle CORS.

		// No more middlewares, pass to the chain.
		web.FinalHandler(router),
	)

}
