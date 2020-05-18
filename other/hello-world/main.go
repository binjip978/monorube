package main

import (
	"net/http"
	"log"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}