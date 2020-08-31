package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"math/rand"
	"os"
	"time"

	"github.com/binjip978/monorube/grpc-k8s-example/server/service"
)

type cfg struct {
	serverAddr string
}

func parseCfg() cfg {
	var serverAddr = flag.String("serverAddr", ":9000", "address of gRPC server endpoint")
	flag.Parse()
	return cfg{serverAddr: *serverAddr}
}

func getID() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown-hostname"
	}
	return hostname
}

func main() {
	cfg := parseCfg()
	fmt.Printf("send data to the server in a loop via gRPC, with current config %+v\n", cfg)
	conn, err := grpc.Dial(cfg.serverAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := service.NewMonitorClient(conn)
	id := getID()
	for {
		client.Point(context.Background(), &service.PointRequest{
			Id: id,
			X:  rand.Int31(),
			Y:  rand.Int31(),
			Z:  rand.Int31(),
		})
		time.Sleep(10 * time.Second)
	}
}
