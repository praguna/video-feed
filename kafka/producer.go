package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var producer *kafka.Producer

func InitProducer() {

	broker := "localhost:9092"

	var err error
	producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

	if err != nil {
		log.Print("Failed to create producer")
		panic(err)
	}

	log.Printf("Created Producer %v", producer)

}

func Produce(topic string, value string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Print("Producing ", topic, " as ", value)
	deliveryChan := make(chan kafka.Event)

	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(value),
		Headers:        []kafka.Header{{Key: "InboundTopic", Value: []byte("InboundTopic feed header value")}},
	}, deliveryChan)

	e := <-deliveryChan

	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("Delivery failed: %v", m.TopicPartition.Error)
	} else {
		log.Print("Initializing Kafka")
		log.Printf("Produced %s [%d] at offset %v",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	if err != nil {
		log.Printf("Error in writing value : %v ", err)
	}
	close(deliveryChan)
}
