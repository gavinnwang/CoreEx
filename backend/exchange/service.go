package exchange

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/wry-0313/exchange/orderbook"
	"github/wry-0313/exchange/pkg/validator"
	"github/wry-0313/exchange/user"
	"log"
	"sync"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	kafkaTopic   = "orders"
	numConsumers = 3
)

var (
	ErrInvalidSymbol = errors.New("Symbol not found")
)

type Service interface {
	PlaceOrder(input PlaceOrderInput) error
	GetMarketPrice(symbol string) (decimal.Decimal, error)
	StartConsumers(brokerList []string)
	ShutdownConsumers()
}

type service struct {
	validator  validator.Validate
	orderBooks map[string]*orderbook.OrderBook
	producer   sarama.SyncProducer
	Shutdown   chan struct{}
}

func NewService(userRepo user.Repository, validator validator.Validate, brokerList []string) Service {
	producer, err := newProducer(brokerList)
	if err != nil {
		log.Fatalf("Could not create producer: %v", err)
	}

	// set up AAPl orderbok
	orderBooks := make(map[string]*orderbook.OrderBook)
	orderBooks["AAPL"] = orderbook.NewOrderBook("AAPL")

	return &service{
		validator:  validator,
		orderBooks: orderBooks,
		producer:   producer,
		Shutdown:   make(chan struct{}),
	}
}

func (s *service) PlaceOrder(input PlaceOrderInput) error {
	if err := s.validator.Struct(input); err != nil {
		return fmt.Errorf("service: validation error: %w", err)
	}

	// Check the validity of the input symbol
	_, ok := s.orderBooks[input.Symbol]
	if !ok {
		return ErrInvalidSymbol
	}


	inputJSON, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("Failed to serialize order to JSON: %w", err)
	}

	// Produce the serialized Order object to Kafka
	msg := &sarama.ProducerMessage{
		Topic: kafkaTopic,
		Value: sarama.ByteEncoder(inputJSON),
	}

	_, _, err = s.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("Failed to send message: %w", err)
	}

	return nil
}

func (s *service) GetMarketPrice(symbol string) (decimal.Decimal, error) {
	ob, ok := s.orderBooks[symbol]
	if !ok {
		return decimal.Decimal{}, ErrInvalidSymbol
	}
	return ob.MarketPrice(), nil
}

func (s *service) StartConsumers(brokerList []string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Fatal("Failed to start consumer:", err)
	}

	log.Printf("Starting Kafka consumers at offest: %v", sarama.OffsetNewest)

	pc, err := consumer.ConsumePartition(kafkaTopic, 0, sarama.OffsetNewest)
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

	for w := 1; w <= numConsumers; w++ {
		wg.Add(1)
		log.Println("Starting consumer", w)
		go s.consumer(pc, &wg)
	}

	wg.Wait()
}

func (s *service) consumer(pc sarama.PartitionConsumer, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case msg, ok := <-pc.Messages():
			if !ok {
				log.Println("Error consuming message")
				continue
			}
			var order PlaceOrderInput
			err := json.Unmarshal(msg.Value, &order)
			if err != nil {
				log.Println("Failed to deserialize order:", err)
				continue
			}
			userID, error := uuid.Parse(order.UserID)
			if error != nil {
				log.Println("Failed to parse UUID:", error)
				continue
			}
			side, error := orderbook.SideFromString(order.OrderSide)
			if error != nil {
				log.Println("Failed to parse side:", error)
				continue
			}
			log.Printf("Consumer processing: %v\n", order)
			switch order.OrderType {
			case "limit":
				_, err := s.orderBooks[order.Symbol].PlaceLimitOrder(side, userID, decimal.NewFromFloat(order.Volume), decimal.NewFromFloat(order.Price))
				if err != nil {
					log.Println(err)
					continue
				}
			case "market":
				_, err := s.orderBooks[order.Symbol].PlaceMarketOrder(side, userID, decimal.NewFromFloat(order.Volume))
				if err != nil {
					log.Println(err)
					continue
				}
			default:
				log.Println("Invalid order type")
			}
		case err := <-pc.Errors():
			log.Println("Error consuming message: ", err)
		case <-s.Shutdown:
			log.Println("Shutting down consumer")
			return
		}
	}
}

func (s *service) ShutdownConsumers() {
	close(s.Shutdown)
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
// 				log.Printf("Added market price data to time series: %s\n", p)
// 			}
// 		}
// 	 }
// }

// func (ex *Exchange) FetchAndStoreBestBids() {
// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ticker.C:
// 			log.Printf("best bids: %s\n", ex.orderBook.BestBid())
// 		}
// 	}
// }

// func (ex *Exchange) FetchAndStoreBestAsks() {
// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ticker.C:
// 			log.Printf("best asks: %s\n", ex.orderBook.BestAsk())
// 		}
// 	}
// }
