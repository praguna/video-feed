package main

import (
	"fmt"
	"net/http"
	"strconv"
	"video-feed/kafka"
	"video-feed/redis"
	"github.com/gorilla/mux"
)

func VideoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), 405)
		return
	}
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	vd, err := redis.VideoDisplay(id)
	if err == redis.VideoNoError {
		http.Error(w, http.StatusText(400), 400)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	_, err = fmt.Fprintf(w, "Id->%s\n"+
		"Title->%s\n"+
		"category->%s\n"+
		"likes->%d\n", id, vd.Title, vd.Category, vd.Likes)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), 405)
		return
	}
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	go kafka.Produce("likes", id)
}

func PopularHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), 405)
		return
	}
	x := mux.Vars(r)["num"]
	if x == "" {
		x = "3"
	}
	xInt, err := strconv.Atoi(x)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	vds, err := redis.GetPopular(xInt)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	c := 0
	for i, vd := range vds {
		if vd == nil {
			break
		}
		_, err = fmt.Fprintf(w, "No->%d\n"+
			"Title->%s\n"+
			"category->%s\n"+
			"likes->%d\n\n", i, vd.Title, vd.Category, vd.Likes)
		c = c + 1
	}
	if c == 0 {
		fmt.Fprintln(w, "No video liked yet")
	}
}

func KafkaRedisServiceTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), 405)
		return
	}
	w.Write([]byte("KafkaRedisServiceTestHandler"))
}
