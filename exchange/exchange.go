package exchange

import (
	"encoding/json"
	"fmt"
	"github/wry-0313/exchange/orderbook"
	ws "github/wry-0313/exchange/websocket"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

type Exchange struct {
	orderBook *orderbook.OrderBook
	producer  sarama.SyncProducer
	Shutdown  chan struct{}
}

const numWorkers = 4
const topic = "orders"

func NewExchange() *Exchange {
	producer, err := newProducer([]string{"localhost:9092"})
	if err != nil {
		log.Fatalf("Could not create producer: %v", err)
	}

	return &Exchange{
		orderBook: orderbook.NewOrderBook(),
		producer:  producer,
		Shutdown:  make(chan struct{}),
	}
}

func (ex *Exchange) Run() {
	http.HandleFunc("/order", ex.PlaceOrderHandler())
	http.HandleFunc("/price", ex.HandleStreamMarketPrice)

	go ex.RunConsumer([]string{"localhost:9092"})

	go ex.FetchAndStoreBestBids()
	go ex.FetchAndStoreBestAsks()
}

func (ex *Exchange) PlaceOrderHandler() http.HandlerFunc {
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

		orderJSON, err := json.Marshal(order)
		if err != nil {
			fmt.Println("Failed to serialize order:", err)
			return
		}
		// Produce the serialized Order object to Kafka
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(orderJSON), // Use ByteEncoder to send it as bytes
		}

		_, _, err = ex.producer.SendMessage(msg)
		if err != nil {
			fmt.Println("Failed to send message:", err)
		}

		fmt.Fprintf(w, "Order received: %+v", order)
	}
}

func (ex *Exchange) RunConsumer(brokerList []string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Fatal("Failed to start consumer:", err)
	}

	pc, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatal("Failed to start partition consumer:", err)
	}

	defer func() {
		log.Println("Closing consumer")
		if err := pc.Close(); err != nil {	
			panic(err)
		}
	}()
	
	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		log.Println("Starting worker", w)
		go ex.worker(pc, &wg)
	}

	wg.Wait()
}

func (ex *Exchange) worker(pc sarama.PartitionConsumer, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case msg, ok := <-pc.Messages():
			if !ok {
				log.Println("Error consuming message")
				continue
			}
			var order OrderRequestParameter
			err := json.Unmarshal(msg.Value, &order)
			if err != nil {
				fmt.Println("Failed to deserialize order:", err)
				continue
			}
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
				log.Println("Placing limit order")
				_, err := ex.orderBook.PlaceLimitOrder(side, clientID, decimal.NewFromFloat(order.Volume), decimal.NewFromFloat(order.Price))
				if err != nil {
					fmt.Println(err)
					continue
				}
			case "market":
				log.Println("Placing market order")
				_, err := ex.orderBook.PlaceMarketOrder(side, clientID, decimal.NewFromFloat(order.Volume))
				if err != nil {
					fmt.Println(err)
					continue
				}
			default:
				fmt.Println("Invalid order type")
			}
		case err := <-pc.Errors():
			log.Println("Error consuming message: ", err)
		case <-ex.Shutdown: 
			log.Println("Shutting down worker")
			return
		}
	}
}

// func (ex *Exchange) FetchAndStoreMarketPrice() {
// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ticker.C:
// 			p := ex.OrderBook.MarketPrice()
// 			timestamp := time.Now().Unix()
// 			key := "market_price"
// 			_, err := ex.rdb.Do("TS.ADD", key, timestamp, p.String()).Result()
// 			if err != nil {
// 				log.Fatalf("Could not add data to time series: %v", err)
// 			} else {
// 				fmt.Printf("Added market price data to time series: %s\n", p)
// 			}
// 		}
// 	}
// }

func (ex *Exchange) FetchAndStoreBestBids() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("best bids: %s\n", ex.orderBook.BestBid())
		}
	}
}

func (ex *Exchange) FetchAndStoreBestAsks() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("best asks: %s\n", ex.orderBook.BestAsk())
		}
	}
}
func (ex *Exchange) HandleStreamMarketPrice(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p := ex.orderBook.MarketPrice()
			priceString := p.String()

			if err := conn.WriteMessage(websocket.TextMessage, []byte(priceString)); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
