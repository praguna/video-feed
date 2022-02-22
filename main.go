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
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong !!\n"))
	})

	router.HandleFunc("/produce-to-incoming-topic", ProduceToIncomingTopic)

	router.HandleFunc("/consume-from-outgoing-topic", ConsumeFromOutgoingTopic)

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
