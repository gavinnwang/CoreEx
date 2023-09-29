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
	"time"

	"github.com/IBM/sarama"
	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
)

const (
	kafkaTopic    = "orders"
	NumPartitions = 5
)

var (
	ErrInvalidSymbol = errors.New("Symbol not found")
)

type Service interface {
	PlaceOrder(input PlaceOrderInput) error

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
		ID:    marketSimulationUlid.String(),
		Name:  "Market Simulation",
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
		ob.SimulateMarketFluctuations(marketSimulationUlid)
		time.Sleep(4 * time.Second)
		ob.Run()
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

func (s *service) startConsumers(brokerList []string) {
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin(brokerList, config)
	if err != nil {
		log.Fatal("Error while creating cluster admin: ", err.Error())
	}
	defer func() { admin.Close() }()
	topics, err := admin.ListTopics()
	if err != nil {
		log.Fatal("Error listing topics: ", err.Error())
	}

	if _, exists := topics[kafkaTopic]; !exists {
		err = admin.CreateTopic(kafkaTopic, &sarama.TopicDetail{
			NumPartitions:     NumPartitions,
			ReplicationFactor: 1,
		}, false)
		if err != nil {
			log.Fatal("Error while creating topic: ", err.Error())
		}
	} else {
		log.Printf("Topic '%s' already exists. Skipping creation.", kafkaTopic)
	}

	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Fatal("Failed to start consumer:", err)
	}

	var wg sync.WaitGroup

	partitionList, err := consumer.Partitions(kafkaTopic)
	if err != nil {
		log.Fatal("Failed to get the list of partitions:", err)
	}

	for _, partition := range partitionList {

		pc, err := consumer.ConsumePartition(kafkaTopic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Failed to start consumer for partition %d: %s", partition, err)
		}

		go func(pc sarama.PartitionConsumer) {
			<-s.Shutdown
			log.Printf("Shutting down partition consumer %d\n", partition)
			pc.AsyncClose()
		}(pc)

		go s.consume(pc, &wg, 0)
	}

	wg.Wait()
}

func (s *service) consume(pc sarama.PartitionConsumer, wg *sync.WaitGroup, index int) {
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
			log.Printf("Shutting down partition consumer %d\n", index)
			return
		}
	}
}

func (s *service) ShutdownConsumers() {
	log.Println("Shutting down consumers called")
	close(s.Shutdown)
}
