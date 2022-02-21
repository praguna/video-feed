package main

import (
	"fmt"
	"net/http"
	"video-feed/kafka"
)

func ProduceToOutgoingTopic(w http.ResponseWriter, r *http.Request) {

	message, ok := r.URL.Query()["message"]
	if !ok || len(message[0]) < 1 {
		fmt.Println("Url Param 'message' is missing")
		return
	}

	go kafka.Produce("messages", message[0])
}

// func ConsumeFromIncomingTopic(w http.ResponseWriter, r *http.Request){

// }
