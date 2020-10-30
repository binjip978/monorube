package main

import (
	"context"
	"log"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

func main() {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()
	data := []byte(`{"hello": "world"}`)
	ctx := context.Background()
	for {
		err = client.PublishEvent(ctx, "pubsub", "test-topic", data)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(60 * time.Second)
	}
}
