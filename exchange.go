package main

import (
	"encoding/json"
	"fmt"
	"github/wry-0313/exchange/orderbook"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Exchange struct {
	OrderBook *orderbook.OrderBook
	rdb       *redis.Client
}

func NewExchange() *Exchange {
	return &Exchange{
		OrderBook: orderbook.NewOrderBook(),
		rdb:       NewRedis(),
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

const numWorkers = 4

func worker(ex *Exchange, jobs <-chan *sarama.ConsumerMessage) {
	for msg := range jobs {
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
	}
}

func (ex *Exchange) consumeAndPlaceOrders(brokerList []string) {
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

	jobs := make(chan *sarama.ConsumerMessage, 100)

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		go worker(ex, jobs)
	}

	for {
		select {
		case msg := <-pc.Messages():
			jobs <- msg
		case err := <-pc.Errors():
			log.Println("Error consuming message: ", err)
		}
	}
}

func (ex *Exchange) FetchMarketPrice() {
	for {
		p := ex.OrderBook.MarketPrice()
		// timestamp := time.Now().Unix()

		// key := "market_price"
		// rdb.Add(ctx, key, []redis.Z{{Score: float64(timestamp), Member: fakePrice}})

		fmt.Printf("market price: %s\n", p)
		time.Sleep(1 * time.Second)
	}
}

func (ex *Exchange) FetchBestBids() {
	for {
		fmt.Printf("best bids: %s\n", ex.OrderBook.BestBid())
		time.Sleep(1 * time.Second)
	}
}

func (ex *Exchange) FetchBestAsks() {
	for {
		fmt.Printf("best asks: %s\n", ex.OrderBook.BestAsk())
		time.Sleep(1 * time.Second)
	}
}
