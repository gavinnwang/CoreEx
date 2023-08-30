package main

import (
	"log"
	"net/http"
)


func main() {
	// Kafka configuration
	brokerList := []string{"localhost:9092"}
	producer, err := newProducer(brokerList)
	if err != nil {
		log.Fatalf("Could not create producer: %v", err)
	}

	ex := NewExchange()

	// Start the Kafka consumer in a new goroutine.
	go ex.consumeOrders(brokerList)

	// Set up the HTTP server.
	http.HandleFunc("/order", ex.PlaceOrderHandler(producer))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
