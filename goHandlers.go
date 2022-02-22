package main

import (
	"fmt"
	"net/http"
	"video-feed/kafka"
)

func ProduceToIncomingTopic(w http.ResponseWriter, r *http.Request) {

	message, ok := r.URL.Query()["message"]
	if !ok || len(message[0]) < 1 {
		fmt.Println("Url Param 'message' is missing")
		return
	}

	go kafka.Produce("messages", message[0])
}

func ConsumeFromOutgoingTopic(w http.ResponseWriter, r *http.Request) {
	result, err := kafka.ReadOutgoingFromRedis()
	if err != nil {
		fmt.Println("error")
	}
	w.Write([]byte(result))
}

func StartConsumer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bananas"))
}
