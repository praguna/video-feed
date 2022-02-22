package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func SetRedisTopic(topic string, message string) error {
	fmt.Println("Setting InboundTopic in redis")
	conn := Pool.Get()
	defer conn.Close()
	key := topic
	value := message
	_, err := redis.String(conn.Do("SET", key, value, "EX", 6000))
	if err != nil {
		return err
	}
	return nil
}

func GetRedisTopic(topic string) (string, error) {
	fmt.Println("Getting Redis Topic", topic)

	conn := Pool.Get()
	defer conn.Close()
	message, err := redis.String(conn.Do("GET", topic))
	if err != nil {
		return "BLOOB", err
	}
	fmt.Println(message, "message")
	return message, nil
}
