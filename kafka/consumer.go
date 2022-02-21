package kafka

import "C"
import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"video-feed/redis"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Consumer(topics []string) {
	group := "messages"
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	broker := "localhost:9092,localhost:9093"
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               broker,
		"group.id":                        group,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
	})
	if err != nil {
		log.Panicf("Failed to create consumer: %s\n", err)
	}
	log.Printf("Created Consumer %v\n", c)
	err = c.SubscribeTopics(topics, nil)
	for {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			_ = c.Close()
			os.Exit(1)

		case ev := <-c.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				c.Unassign()
			case *kafka.Message:
				handleRedis(*e.TopicPartition.Topic, string(e.Value))
			case kafka.PartitionEOF:
				log.Printf("%% Reached %v\n", e)
			case kafka.Error:
				log.Printf("%% Error: %v\n", e)
			}
		}
	}
}

func handleRedis(topic string, value string) {
	fmt.Println("\n\nInitializing Redis\n\n\n")
	switch topic {
	case "messages":
		err := redis.SetMessage(value)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("\n\n\nSet %v topic in redis to '%v'\n\n ********* \n\nEND\n\n *********\n\n", topic, value)
		}
	}
}
