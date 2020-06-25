package main

import (
	"log"
	"net/http"
	"os"
	"video-feed/kafka"
	"video-feed/redis"
)

func main() {
	if len(os.Args) > 1{
		if os.Args[1] == "populate"{
			redis.Populate()
			return
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong !!\n"))
	})

	mux.HandleFunc("/video", VideoHandler)

	mux.HandleFunc("/like", LikeHandler)

	mux.HandleFunc("/popular", PopularHandler)

	mux.HandleFunc("/upload",func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(200)
		return
	})
	log.Println("Initializing redis pool: ")
	redis.Init()
	go kafka.Consumer([]string{"likes", "upload"})
	log.Println("Video-Feed Listening on :4000")
	err := http.ListenAndServe(":4000", mux)
	if err == nil {
		log.Printf("Server error %v :", err)
	}
}
