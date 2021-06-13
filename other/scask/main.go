package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type cfg struct {
	port int
}

func parseCfg() cfg {
	c := cfg{
		port: *flag.Int("port", 8011, "server port"),
	}
	flag.Parse()

	return c
}

var errKeyNotFount = errors.New("key not found")

type api interface {
	get(key string) (string, error)
	put(key string, value string) error
	del(key string) error
}

type memStore struct {
	sync.RWMutex
	storage map[string]string
}

func newMemStore() *memStore {
	return &memStore{storage: make(map[string]string)}
}

func (mem *memStore) get(key string) (string, error) {
	mem.RLock()
	defer mem.RUnlock()
	value, ok := mem.storage[key]
	if !ok {
		return "", errKeyNotFount
	}

	return value, nil
}

func (mem *memStore) put(key string, value string) error {
	mem.Lock()
	defer mem.Unlock()
	mem.storage[key] = value

	return nil
}

func (mem *memStore) del(key string) error {
	mem.Lock()
	defer mem.Unlock()
	delete(mem.storage, key)

	return nil
}

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
	if r.Method != http.MethodGet {
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

func main() {
	cfg := parseCfg()
	log.Printf("config: %+v\n", cfg)

	backend := newMemStore()
	srv := newServer(cfg, backend)

	log.Fatal(srv.httpServer.ListenAndServe())
}
