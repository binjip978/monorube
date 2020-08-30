package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"server/service"
)

type cfg struct {
	serverAddr string
}

func parseConfig() cfg {
	var serverAddr = flag.String("serverAddr", ":9000", "grpc server host:port")
	flag.Parse()
	return cfg{serverAddr: *serverAddr}
}

func main() {
	cfg := parseConfig()
	fmt.Printf("starting gRPC server with: %+v config\n", cfg)
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	service.RegisterMonitorServer(grpcServer, &service.PointServer{})
	log.Fatal(grpcServer.Serve(lis))
}
