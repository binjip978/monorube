package main

import (
	"io/ioutil"
	"log"
	"tftp/tftp"
)

func main() {
	p, err := ioutil.ReadFile("./1.jpg")
	if err != nil {
		panic(err)
	}

	s := tftp.Server{Payload: p}
	log.Fatal(s.ListenAndServe("127.0.0.1:6900"))
}
