package main

import (
	"bufio"
	"kq-echo/kqueue"
	"kq-echo/socket"
	"log"
	"strings"
)

func main() {
	s, err := socket.Listen("127.0.0.1", 8080)
	if err != nil {
		panic(err)
	}

	loop, err := kqueue.NewEventLoop(s)
	if err != nil {
		panic(err)
	}

	log.Println("server started on 8080")

	loop.Handle(func(s *socket.Socket) {
		log.Println("handle connction")
		reader := bufio.NewReader(s)
		for {
			l, err := reader.ReadString('\n')
			if err != nil || strings.TrimSpace(l) == "" {
				break
			}
			s.Write([]byte(l))
		}
		s.Close()
	})
}
