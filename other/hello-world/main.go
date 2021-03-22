package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World Old"))
	})

	http.HandleFunc("/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Println("sleeping for a 2 seconds")
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("v1 Hello World"))
	})

	http.HandleFunc("/v2/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Println("sleeping for a 2 minutes")
		time.Sleep(2 * time.Minute)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("v2 Hello World Sleepy"))
	})

	http.HandleFunc("/v3/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Println("fast version, no sleep at all")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("v3 Hello World Fast"))
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}
