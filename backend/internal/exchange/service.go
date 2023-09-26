package exchange

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/wry-0313/exchange/internal/models"
	"github/wry-0313/exchange/internal/orderbook"
	"github/wry-0313/exchange/internal/user"
	"github/wry-0313/exchange/pkg/validator"
	"log"
	"sync"

	"github.com/IBM/sarama"
	"github.com/oklog/ulid/v2"
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
	GetMarketPrice(symbol string) (float64, error)
	GetSymbolInfo(symbol string) (orderbook.SymbolInfoResponse, error)
	Run(brokerList []string)
	ShutdownConsumers()
	GetSymbolMarketPriceHistory(symbol string) ([]models.StockPriceHistory, error)
}

type service struct {
	validator  validator.Validate
	obServices map[string]orderbook.Service
	producer   sarama.SyncProducer
	Shutdown   chan struct{}
	userRepo   user.Repository
}

func NewService(userRepo user.Repository, obServices map[string]orderbook.Service, validator validator.Validate, brokerList []string) Service {
	producer, err := newProducer(brokerList)
	if err != nil {
		log.Fatalf("Could not create producer: %v", err)
	}
	return &service{
		validator:  validator,
		obServices: obServices,
		producer:   producer,
		Shutdown:   make(chan struct{}),
		userRepo:   userRepo,
	}
}

func (s *service) Run(brokerList []string) {
	marketSimulationUlid := ulid.Make()
	email := "market@gmail.com"
	err := s.userRepo.CreateUser(models.User{
		ID: marketSimulationUlid.String(),
		Name: "Market Simulation",
		Email: &email,
	})

	if err != nil {
		// log.Fatalf("Service: failed to create market simulation user: %v", err)
		// retrieve the ulid if already exist
		ulidString, err := s.userRepo.GetUserByEmail(email)
		if err != nil {
			log.Fatalf("Service: failed to retrieve market simulation user: %v", err)
		}
		marketSimulationUlid, err = ulid.Parse(ulidString.ID)
		if err != nil {
			log.Fatalf("Service: failed to parse market simulation user ulid: %v", err)
		}
	}

	go s.startConsumers(brokerList)
	for _, ob := range s.obServices {
		log.Printf("Starting market price history persistance for %v\n", ob.Symbol())
		ob.Run()
		ob.SimulateMarketFluctuations(marketSimulationUlid)
	}
}

func (s *service) GetSymbolMarketPriceHistory(symbol string) ([]models.StockPriceHistory, error) {
	ob, ok := s.obServices[symbol]
	if !ok {
		return nil, ErrInvalidSymbol
	}
	return ob.GetMarketPriceHistory()
}

func (s *service) PlaceOrder(input PlaceOrderInput) error {
	if err := s.validator.Struct(input); err != nil {
		return fmt.Errorf("service: validation error: %w", err)
	}

	// Check the validity of the input symbol
	_, ok := s.obServices[input.Symbol]
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

func (s *service) GetMarketPrice(symbol string) (float64, error) {
	ob, ok := s.obServices[symbol]
	if !ok {
		return 0, ErrInvalidSymbol
	}
	return ob.MarketPrice().InexactFloat64(), nil
}

func (s *service) GetSymbolInfo(symbol string) (orderbook.SymbolInfoResponse, error) {
	ob, ok := s.obServices[symbol]
	if !ok {
		return orderbook.SymbolInfoResponse{}, ErrInvalidSymbol
	}
	return orderbook.SymbolInfoResponse{
		Symbol:    symbol,
		AskVolume: ob.AskVolume().InexactFloat64(),
		BidVolume: ob.BidVolume().InexactFloat64(),
		BestAsk:   ob.BestAsk().InexactFloat64(),
		BestBid:   ob.BestBid().InexactFloat64(),
		Price:     ob.MarketPrice().InexactFloat64(),
	}, nil

}

func (s *service) startConsumers(brokerList []string) {
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
		log.Println("Closing consumer partition")
		if err := pc.Close(); err != nil {
			panic(err)
		}
	}()

	var wg sync.WaitGroup

	for w := 1; w <= numConsumers; w++ {
		wg.Add(1)
		log.Println("Starting consumer", w)
		go s.consumer(pc, &wg, w)
	}

	wg.Wait()
}

func (s *service) consumer(pc sarama.PartitionConsumer, wg *sync.WaitGroup, index int) {
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
			userID, error := ulid.Parse(order.UserID)
			if error != nil {
				log.Println("Failed to parse ULID:", error)
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
				_, err := s.obServices[order.Symbol].PlaceLimitOrder(side, userID, decimal.NewFromFloat(order.Volume).Round(2), decimal.NewFromFloat(order.Price).Round(2))
				if err != nil {
					log.Println(err)
					continue
				}
			case "market":
				_, err := s.obServices[order.Symbol].PlaceMarketOrder(side, userID, decimal.NewFromFloat(order.Volume).Round(2))
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
			log.Printf("Shutting down consumer %d\n", index)
			return
		}
	}
}

func (s *service) ShutdownConsumers() {
	log.Println("Shutting down consumers called")
	close(s.Shutdown)
}
