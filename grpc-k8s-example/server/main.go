package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/binjip978/monorube/grpc-k8s-example/server/service"
	"google.golang.org/grpc"
)

var hostname string

type cfg struct {
	serverAddr     string
	httpServerAddr string
}

func parseConfig() cfg {
	var serverAddr = flag.String("serverAddr", ":9000", "gRPC server host:port")
	var httpServerAddr = flag.String("httpServerAddr", ":8000", "http server host:port")
	flag.Parse()
	return cfg{serverAddr: *serverAddr, httpServerAddr: *httpServerAddr}
}

func httpServer(serverAddr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s-%s", hostname, "v0.0.1")))
	})

	srv := &http.Server{
		Addr:         serverAddr,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Handler:      mux,
	}

	return srv
}

func main() {
	cfg := parseConfig()
	fmt.Printf("starting gRPC server with: %+v config\n", cfg)
	hostname, _ = os.Hostname()
	srv := httpServer(cfg.httpServerAddr)
	go func() {
		log.Println(srv.ListenAndServe())
	}()

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	service.RegisterMonitorServer(grpcServer, &service.PointServer{})
	log.Fatal(grpcServer.Serve(lis))
}
