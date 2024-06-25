package main

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
)

var testDBPath string = "test.db"

func TestMain(m *testing.M) {
	_, err := exec.LookPath("sqlite3")
	if err != nil {
		panic("sqlite3 binary not found.")
	}

	err = exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("sqlite3 %s \"create table auth(ID INTEGER, name, pass TEXT, PRIMARY KEY(ID)); insert into auth values(1, 'admin', 'admin');\"", testDBPath),
	).Run()
	if err != nil {
		os.Remove(testDBPath)
		panic("error setting up test db: " + err.Error())
	}

	status := m.Run()
	os.Remove(testDBPath)
	os.Exit(status)
}

func TestServer(t *testing.T) {
	router, err := setupRouter(testDBPath)
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
