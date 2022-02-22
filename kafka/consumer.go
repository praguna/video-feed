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
	group := "InboundTopic"
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
				fmt.Println("In consumer")
				myTopic := *e.TopicPartition.Topic
				fmt.Printf("myTopic Type = %T", myTopic)
				if myTopic == "InboundTopic" {
					saveRedisTriggerOutboundTopicKafka(*e.TopicPartition.Topic, string(e.Value))
				} else {
					emitToRedis(*e.TopicPartition.Topic, string(e.Value))
				}
			case kafka.PartitionEOF:
				log.Printf("%% Reached %v\n", e)
			case kafka.Error:
				log.Printf("%% Error: %v\n", e)
			}
		}
	}
}

func emitToRedis(topic string, value string) {
	fmt.Println("\n\nEmitting to Redis\n\n\n")
	fmt.Println(("IN InboundTopic"))
	err := redis.SetRedisTopic(topic, value)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("\n\n\nSet %v topic in redis to '%v'\n\n ********* \n\nEND\n\n *********\n\n", topic, value)
	}
}
func readFromRedis(topic string) (string, error) {
	fmt.Println("\n\nReading from Redis\n\n\n")
	fmt.Println("topic", topic)
	result, err := redis.GetRedisTopic(topic)
	if err != nil {
		return "", err
	}
	fmt.Printf("result %s", result)
	return result, nil
}
func reverseString(str string) (string, error) {
	rune_arr := []rune(str)
	var rev []rune
	for i := len(rune_arr) - 1; i >= 0; i-- {
		rev = append(rev, rune_arr[i])
	}
	fmt.Println("Reverse: ", string(rev))
	result := string(rev)
	return result, nil
}

func produceOutgoingTopic(str string) {
	Produce("OutgoingTopic", str)
}
func saveRedisTriggerOutboundTopicKafka(topic string, value string) error {
	fmt.Println("Running Sequence")
	// save topic to redis
	emitToRedis(topic, value)

	// read topic from redis
	message, err := readFromRedis(topic)
	if err != nil {
		return err
	}
	fmt.Println("Received Message: ", message)

	// reverse message
	reversed, err := reverseString(message)
	if err != nil {
		return err
	}

	// produce reversed message to Kafka
	produceOutgoingTopic(reversed)
	return nil
}
