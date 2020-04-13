package main

import (
	"flag"
	"io"
	"log"
	"net"
)

type args struct {
	proxyAddr string
	dstAddr   string
}

func parseArgs() args {
	proxyAddr := flag.String("proxyAddr", ":8080", "proxy addr")
	dstAddr := flag.String("dstAddr", "localhost:9080", "dst addr")
	flag.Parse()

	return args{*proxyAddr, *dstAddr}
}

func main() {
	args := parseArgs()
	lis, err := net.Listen("tcp", args.proxyAddr)
	if err != nil {
		log.Fatalf("can't bind a port: %v", err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("can't accept connection: %v", err)
		}
		go handle(conn, args.dstAddr)
	}
}

func handle(src net.Conn, dstAddr string) {
	log.Printf("incomming connetion: %s", src.RemoteAddr().String())
	dst, err := net.Dial("tcp", dstAddr)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer dst.Close()

	go func() {
		if _, err := io.Copy(dst, src); err != nil {
			log.Fatalf("can' copy from dst -> src: %v", err)
		}
	}()

	if _, err := io.Copy(src, dst); err != nil {
		log.Fatalf("can' copy from src -> dst: %v", err)
	}
}
