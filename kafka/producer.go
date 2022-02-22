package kafka

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var producer *kafka.Producer

func InitProducer() {

	broker := "localhost:9092"

	var err error
	producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

	if err != nil {
		log.Println("Failed to create producer")
		panic(err)
	}

	log.Printf("Created Producer %v\n", producer)

}

func Produce(topic string, value string) {
	log.Printf("\n\n\n ********* \n\nSTART\n\n *********\n\n\n\n")

	fmt.Println("Producing", topic, "as", value)
	deliveryChan := make(chan kafka.Event)

	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(value),
		Headers:        []kafka.Header{{Key: "InboundTopic", Value: []byte("InboundTopic feed header value")}},
	}, deliveryChan)

	e := <-deliveryChan

	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		log.Printf("\n\n\n\nInitializing Kafka\n\n\n")
		log.Printf("\n\n\n\nProduced %s [%d] at offset %v\n\n\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	if err != nil {
		log.Printf("Error in writing value : %v \n", err)
	}
	close(deliveryChan)
}
