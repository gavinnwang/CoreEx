package main

import (
	"encoding/json"
	"fmt"
	"github/wry-0313/exchange/orderbook"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Exchange struct {
	OrderBook *orderbook.OrderBook
}

func NewExchange() *Exchange {
	return &Exchange{
		OrderBook: orderbook.NewOrderBook(),
	}
}

func (ex *Exchange) PlaceOrderHandler(producer sarama.SyncProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		// Decode the JSON request into an Order object
		var order OrderRequestParameter
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&order)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// // Serialize the Order object to JSON to produce it to Kafka
		// orderJSON, err := json.Marshal(order)
		if err != nil {
			http.Error(w, "Failed to serialize order to JSON", http.StatusInternalServerError)
			return
		}

		// Produce the serialized Order object to Kafka
		produce(producer, "orders", &order)

		fmt.Fprintf(w, "Order received: %+v", order)
	}
}


// consumeOrders reads orders from a Kafka topic and processes them.
func (ex *Exchange) consumeOrders(brokerList []string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start consumer:", err)
	}

	pc, err := consumer.ConsumePartition("orders", 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalln("Failed to start partition consumer:", err)
	}

	defer pc.Close()

	for {
		select {
		case msg := <-pc.Messages():
			var order OrderRequestParameter
			err := json.Unmarshal(msg.Value, &order)
			if err != nil {
				fmt.Println("Failed to deserialize order:", err)
				continue
			}
			fmt.Printf("Consumed message: %v\n", order)
			clientID, error := uuid.Parse(order.ClientID)
			if error != nil {
				fmt.Println("Failed to parse UUID:", error)
				continue
			}
			side, error := orderbook.SideFromString(order.OrderSide)
			if error != nil {
				fmt.Println("Failed to parse side:", error)
				continue
			}
			switch order.OrderType {
			case "limit":
				_, err := ex.OrderBook.PlaceLimitOrder(side, clientID, decimal.NewFromFloat(order.Volume), decimal.NewFromFloat(order.Price))
				if err != nil {
					fmt.Println(err)
					continue
				}
			case "market":
				_, err := ex.OrderBook.PlaceMarketOrder(side, clientID, decimal.NewFromFloat(order.Volume))
				if err != nil {
					fmt.Println(err)
					continue
				}
			default:
				fmt.Println("Invalid order type")
			}
		case err := <-pc.Errors():
			log.Println("Error consuming message: ", err)
		}
	}
}