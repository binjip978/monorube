package main

import (
	"log"
	"net/http"
)

func isPrime(number uint64) bool {
	var i uint64 = 2
	for i < number {
		if number%2 == 0 {
			return false
		}
	}

	return true
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz)
	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}
