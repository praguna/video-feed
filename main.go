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

	// TODO: Set this route to start a perpetual consumer that listens for topic updates
	router.HandleFunc("/start-consumer", StartConsumer)

	router.HandleFunc("/consume-from-outgoing-topic", ConsumeFromOutgoingTopic)

	router.HandleFunc("/produce-to-incoming-topic", ProduceToIncomingTopic)

	http.Handle("/", router)

	log.Println("Initializing redis pool: ")
	redis.Init()
	go kafka.InitProducer()
	go kafka.Consumer([]string{"messages", "outgoing"})
	log.Println("Messages Topic App Listening on :4000")
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		log.Printf("Server error %v :", err)
	}
}
