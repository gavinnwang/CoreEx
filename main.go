package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Producer function
func producer() {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	topic := "foobar"

	// Produce messages to topic (asynchronously)
	for _, word := range []string{"Order1", "Order2", "Order3"} {
		producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}

	time.Sleep(5 * time.Second)
	// Produce messages: numbers 1 to 9999
	for i := 1; i <= 9999; i++ {
		value := strconv.Itoa(i) // Convert integer to string

		// Produce message to the Kafka topic
		producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(value),
		}, nil)

		// Poll for events, including delivery reports and errors
	}
	// Wait for delivery reports
	for e := range producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
			} else {
				fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}

}

// Consumer function
func consumer() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	consumer.SubscribeTopics([]string{"foobar"}, nil)
	fmt.Println("topic subscribed")

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received message: %s\n", msg.Value)
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}
}

func main() {
	// Done channels for producer and consumer
	doneProducer := make(chan bool)
	doneConsumer := make(chan bool)

	// Run producer and consumer concurrently
	go func() {
		fmt.Println("producer")
		producer()
		doneProducer <- true
	}()
	time.Sleep(1 * time.Second)

	go func() {
		fmt.Println("consumer")
		consumer()

		doneConsumer <- true
	}()

	// Wait for producer and consumer to finish
	<-doneProducer
	<-doneConsumer
}
