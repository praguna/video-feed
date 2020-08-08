package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

var producer *kafka.Producer

func InitProducer(){
	broker := "localhost:9092"
	var err error
	producer,err =  kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil{
		log.Println("Failed to create producer")
		panic(err)
	}
	log.Printf("Created Producer %v\n", producer)
}

func Produce(topic string, value string){
	deliveryChan := make(chan kafka.Event)
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(value),
		Headers:        []kafka.Header{{Key: "videoFeed", Value: []byte("video feed header value")}},
	}, deliveryChan)

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		log.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	if err!=nil{
		log.Printf("Error in writing value : %v \n",err)
	}
	close(deliveryChan)
}