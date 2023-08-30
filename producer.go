package main

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

// newProducer creates a new Kafka producer.
func newProducer(brokerList []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

// produce sends a message to a Kafka topic.
func produce(producer sarama.SyncProducer, topic string, order *OrderRequestParameter) {
	orderJSON, err := json.Marshal(order)
	if err != nil {
		fmt.Println("Failed to serialize order:", err)
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(orderJSON), // Use ByteEncoder to send it as bytes
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Failed to send message:", err)
	}
}
