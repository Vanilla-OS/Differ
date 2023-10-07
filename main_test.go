package main

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	router, err := setupRouter()
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/status", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("/status request returned non-OK status '%d'", w.Code)
	}

	if w.Body.String() != "{\"status\":\"ok\"}" {
		t.Fatalf("/status request returned unexpected body '%s'", w.Body.String())
	}
}
