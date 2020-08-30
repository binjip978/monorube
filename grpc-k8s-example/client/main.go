package main

import (
	"flag"
	"fmt"
)

type cfg struct {
	serverAddr string
}

func parseCfg() cfg {
	var serverAddr = flag.String("serverAddr", ":9000", "address of gRPC server endpoint")
	flag.Parse()
	return cfg{serverAddr: *serverAddr}
}

func main() {
	cfg := parseCfg()
	fmt.Printf("send data to the server in a loop via gRPC, with current config %+v\n", cfg)
}
