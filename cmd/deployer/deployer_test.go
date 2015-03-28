package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCliCreateRelease(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, r.PostFormValue("builds"), "[]")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, fixtures.release1Data)
	}))
	defer ts.Close()

	app := CreateCliApp(ts.URL)
	app.Run([]string{"release", "service1", "tag1", "service2", "tag2"})
	//assert.Contains(recorder.Body, `{json}\nid`)
}

func TestCliCreateReleaseAndDeploy() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(r.PostFormValue("builds"), "[]")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, fixtures.release1Data)
	}))
	defer ts.Close()

	// same as above, just add extra different body field
	app := CreateCli(config, ts.URL, recorder)
	app.Run([]string{"--to", "env0", "release", "service1", "tag1"})
	recorder.should.contain(`{json}\nid`)
}

func TestCliCreateDeploy() {
	// check for src dest
	app := CreateCli(config, ts.URL, recorder)
	app.Run([]string{"deploy", "release0", "env1"})
}
