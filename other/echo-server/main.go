package main

import (
	"io"
	"log"
	"net"
)

func echo(conn net.Conn) {
	defer conn.Close()
	if _, err := io.Copy(conn, conn); err != nil {
		log.Fatalf("unable to echo: %v", err)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("can't bind a port: %v", err)
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("can't accept connection: %v", err)
		}
		go echo(conn)
	}
}
