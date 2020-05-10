package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsEndpoint(t *testing.T) {
	srv := newServer(config{})
	ts := httptest.NewServer(srv.Handler)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/metrics")
	if err != nil {
		t.Errorf("metrics endpoint should be avaliable, got: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("status code should be 200, got: %v", res.StatusCode)
	}
}
