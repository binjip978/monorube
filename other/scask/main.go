package main

import (
	"flag"
	"log"
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

func main() {
	cfg := parseCfg()
	log.Printf("config: %+v\n", cfg)

	backend := newMemStore()
	srv := newServer(cfg, backend)

	log.Fatal(srv.httpServer.ListenAndServe())
}
