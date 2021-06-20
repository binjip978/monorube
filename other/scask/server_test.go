package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerAPI(t *testing.T) {
	memStore := newMemStore()
	srv := newServer(cfg{port: 8011}, memStore)
	ts := httptest.NewServer(srv.httpServer.Handler)
	defer ts.Close()

	// put k1:v1
	req := request{
		Key:   "k1",
		Value: "v1",
	}

	b, _ := json.Marshal(&req)

	r, err := http.Post(ts.URL+"/v1/put", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Error(err)
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		t.Errorf("status should be 200, was %d", r.StatusCode)
	}

	// get k1
	req = request{Key: "k1"}
	b, _ = json.Marshal(&req)

	r, err = http.Post(ts.URL+"/v1/get", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Error(err)
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		t.Errorf("status should be 200, was: %d", r.StatusCode)
	}

	b, err = ioutil.ReadAll(r.Body)
	if err != nil {
		t.Error(err)
	}

	var res response
	err = json.Unmarshal(b, &res)
	if err != nil {
		t.Error(err)
	}

	if res.Error != "" || res.Key != "k1" || res.Value != "v1" {
		t.Error("don't return what expected, nil err, k1, v1")
	}
}

func TestBitcaskInstance(t *testing.T) {
	storageDir := t.TempDir()
	cask, err := newSCask(storageDir)
	if err != nil {
		t.Fatal(err)
	}

	srv := newServer(cfg{port: 8011}, cask)
	ts := httptest.NewServer(srv.httpServer.Handler)
	defer ts.Close()

	numRecords := 10

	for i := 0; i < numRecords; i++ {
		req := request{
			Key:   fmt.Sprintf("k-%d", i),
			Value: fmt.Sprintf("v-%d", i),
		}

		b, _ := json.Marshal(&req)
		r, err := http.Post(ts.URL+"/v1/put", "application/json", bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}
		r.Body.Close()
	}

	for i := 0; i < numRecords; i++ {
		req := request{
			Key: fmt.Sprintf("k-%d", i),
		}

		b, _ := json.Marshal(&req)
		r, err := http.Post(ts.URL+"/v1/get", "application/json", bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}

		b, err = ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		r.Body.Close()

		var res response
		err = json.Unmarshal(b, &res)
		if err != nil {
			t.Fatal(err)
		}

		if res.Value != fmt.Sprintf("v-%d", i) {
			t.Errorf("wrong value, expect v-%d, got %s", i, res.Value)
		}
	}
}
