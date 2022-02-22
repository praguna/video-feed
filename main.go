package main

//////////
import (
	"log"
	"net/http"
	"video-feed/kafka"
	"video-feed/redis"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/produce-to-incoming-topic", ProduceToIncomingTopic)

	http.Handle("/", router)

	log.Println("Initializing redis pool: ")
	redis.Init()
	go kafka.InitProducer()
	go kafka.Consumer([]string{"InboundTopic", "OutgoingTopic"})
	log.Println("InboundTopic Topic App Listening on :4000")
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		log.Printf("Server error %v :", err)
	}
}
