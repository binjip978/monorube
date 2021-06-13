package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

type server struct {
	backend    api
	httpServer *http.Server
}

func (s *server) put(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req request
	err = json.Unmarshal(b, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.backend.put(req.Key, req.Value)
	if err == nil {
		return
	}

	res := response{
		Error: err.Error(),
	}

	b, err = json.Marshal(&res)
	if err != nil {
		panic(err)
	}

	w.Write(b)
}

func (s *server) get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req request
	err = json.Unmarshal(b, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := response{
		Key: req.Key,
	}

	value, err := s.backend.get(req.Key)
	if err != nil {
		res.Error = err.Error()
	} else {
		res.Value = value
	}

	b, err = json.Marshal(&res)
	if err != nil {
		panic(err)
	}

	w.Write(b)
}

func (s *server) del(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req request
	err = json.Unmarshal(b, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.backend.del(req.Key)
	if err == nil {
		return
	}

	res := response{
		Error: err.Error(),
	}

	b, err = json.Marshal(&res)
	if err != nil {
		panic(err)
	}

	w.Write(b)
}

type request struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

type response struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
	Error string `json:"error,omitempty"`
}

func newServer(config cfg, backend api) *server {
	s := &server{
		backend: backend,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/get", s.get)
	mux.HandleFunc("/v1/put", s.put)
	mux.HandleFunc("/v1/del", s.del)

	httpServer := http.Server{
		Addr:    net.JoinHostPort("", strconv.Itoa(config.port)),
		Handler: mux,
	}
	s.httpServer = &httpServer

	return s
}
