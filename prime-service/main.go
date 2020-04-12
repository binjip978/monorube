package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func isPrime(number uint64) bool {
	var i uint64 = 2
	for i < number {
		if number%i == 0 {
			return false
		}

		i++
	}

	return true
}

// Response will return answer to user
type Response struct {
	Number  uint64 `json:"number,omitempty"`
	IsPrime bool   `json:"is_prime"`
	Err     string `json:"err,omitempty"`
}

func (r *Response) bytes() []byte {
	b, _ := json.Marshal(r)
	return b
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func servePrime(w http.ResponseWriter, r *http.Request) {
	numStr := r.URL.Query().Get("n")
	num, err := strconv.ParseUint(numStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		resp := &Response{Err: fmt.Sprintf("can't parse %s as unsigned integer, got %v", numStr, err)}
		w.Write(resp.bytes())
		return
	}

	result := isPrime(num)
	w.WriteHeader(http.StatusOK)
	resp := &Response{Number: num, IsPrime: result}
	w.Write(resp.bytes())
}

func setupServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/prime", servePrime)
	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return &srv
}

func main() {
	srv := setupServer()
	log.Fatal(srv.ListenAndServe())
}
