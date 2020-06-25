package kafka

import "C"
import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"os/signal"
	"syscall"
	"video-feed/redis"
)

func Consumer(topics []string){
	group := "videoFeed"
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
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created Consumer %v\n", c)
	err = c.SubscribeTopics(topics, nil)
	for {
		fmt.Println("kafka consumer active .....")
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			_ = c.Close()
			os.Exit(1)

		case ev := <-c.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				//fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				//fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Unassign()
			case *kafka.Message:
				handleRedis(*e.TopicPartition.Topic, string(e.Value))
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			}
		}
	}
}

func handleRedis(topic string, value string){
	switch topic {
	case "likes":
		err := redis.AddLike(value)
		if err!= nil{
			log.Print(err)
		}else if err == redis.VideoNoError{
			log.Print("Video Id Invalid")
		}
		log.Print("redis updated")
	}
}