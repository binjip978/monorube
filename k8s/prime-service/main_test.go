package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPrime(t *testing.T) {
	var table = []struct {
		number  uint64
		isPrime bool
	}{
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{6, false},
		{7, true},
		{8, false},
		{9, false},
		{10, false},
		{11, true},
	}

	for _, entry := range table {
		if res := isPrime(entry.number); res != entry.isPrime {
			t.Errorf("case for number: %d is not correct, result was %t", entry.number, res)
		}
	}
}

func TestHealthz(t *testing.T) {
	srv := setupServer()
	ts := httptest.NewServer(srv.Handler)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Errorf("get should succeed, got %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("should return 200, got %d", res.StatusCode)
	}
}

func TestServePrime(t *testing.T) {
	srv := setupServer()
	ts := httptest.NewServer(srv.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/prime?n=11")
	if err != nil {
		t.Errorf("get should succeed, got %v", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("should read response body, got %v", err)
	}
	defer resp.Body.Close()
	var response Response
	err = json.Unmarshal(b, &response)
	if err != nil {
		t.Errorf("should parse response body to go struct, got %v", err)
	}
	if response.Number != 11 || !response.IsPrime || response.Err != "" {
		t.Error("response struct is not correct")
	}

	resp, _ = http.Get(ts.URL + "/prime?n=hello")
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("should read response body, got %v", err)
	}
	defer resp.Body.Close()
	err = json.Unmarshal(b, &response)
	if err != nil {
		t.Errorf("should parse response body to go struct, got %v", err)
	}
	if response.Err == "" {
		t.Error("response struct should contain an error")
	}
}

func TestStatus(t *testing.T) {
	srv := setupServer()
	ts := httptest.NewServer(srv.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/status")
	if err != nil {
		t.Errorf("get should succeed, got error: %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("should read response body, got %v", err)
	}

	var st ServiceStatus
	json.Unmarshal(b, &st)
	if err != nil {
		t.Errorf("should parse response body to go struct, got %v", err)
	}

	if st.Hostname == "" {
		t.Error("Hostname should not be empty")
	}
}
