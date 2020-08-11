// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// setupRouting will set all routes for the API.
func (a *API) setupRouting() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Ethereum Swarm resolver")
	})

	a.Handler = router
}
