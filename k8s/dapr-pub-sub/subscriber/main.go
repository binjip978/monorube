package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type subscribe struct {
	PubsubName string `json:"pubsubname"`
	Topic      string `json:"topic"`
	Route      string `json:"route"`
}

func main() {
	http.HandleFunc("/dapr/subscribe", func(w http.ResponseWriter, r *http.Request) {
		log.Println("in subscribe2")
		s := []subscribe{
			{
				PubsubName: "pubsub",
				Topic:      "test-topic",
				Route:      "message",
			},
		}
		b, err := json.Marshal(&s)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	})

	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Panicln(err)
			return
		}
		defer r.Body.Close()
		log.Println(string(b))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
