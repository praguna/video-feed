package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var Pool *redis.Pool

type Video struct{
	Title string `redis:"title"`
	Category string `redis:"category"`
	Likes int `redis:"likes"`
}

func Init(){
	Pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
}