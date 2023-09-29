package test

import (
	"fmt"
	"github/wry-0313/exchange/db"
	"github/wry-0313/exchange/internal/config"
	"github/wry-0313/exchange/internal/orderbook"
	"github/wry-0313/exchange/internal/redis"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"


	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
)

// "fmt"
// "math/rand"
// "reflect"
// "sync"
// "testing"
// "time"

// "github.com/shopspring/decimal"

// func TestPlaceMarketOrderAfterLimit(t *testing.T) {
// 	ob := NewService("AAPL")
// 	clientID := ulid.Make()
// 	for i := 0; i < 2; i++ {
// 		ob.PlaceLimitOrder(orderbook.Buy, clientID, decimal.NewFromInt(10000), decimal.NewFromInt(10))
// 	}
// 	max := ob.BestBid()
// 	fmt.Printf("max: %v\n", max)
// 	ob.PlaceMarketOrder(orderbook.Sell, clientID, decimal.NewFromInt(50))
// 	ob.PlaceMarketOrder(orderbook.Sell, clientID, decimal.NewFromInt(50))
// 	ob.PlaceMarketOrder(orderbook.Sell, clientID, decimal.NewFromInt(50))
// 	ob.PlaceMarketOrder(orderbook.Sell, clientID, decimal.NewFromInt(50))
// 	ob.PlaceMarketOrder(orderbook.Sell, clientID, decimal.NewFromInt(50))
// 	ob.PlaceMarketOrder(orderbook.Sell, clientID, decimal.NewFromInt(9751))
// 	fmt.Printf("max: %v\n", max)
// }

// func TestPlaceMarketOrderAfterLimitConcurrent(t *testing.T) {
// 	ob := NewService("AAPL")
// 	clientID := ulid.Make()
// 	// ch := make(chan int)
// 	var wg sync.WaitGroup
// 	_, err := ob.PlaceLimitOrder(orderbook.Buy, clientID, decimal.NewFromInt(10000), decimal.NewFromInt(10))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	max := ob.BestBid()
// 	fmt.Printf("max: %v\n", max)
// 	time.Sleep(2 * time.Second)
// 	wg.Add(1)
// 	go func() {
// 		fmt.Println("start1")
// 		defer wg.Done()
// 		for i := 0; i < 10; i++ {
// 			ob.PlaceLimitOrder(orderbook.Buy, clientID, decimal.NewFromInt(10), decimal.NewFromInt(10))
// 		}
// 	}()
// 	wg.Add(1)
// 	go func() {
// 		fmt.Println("start1.5")
// 		defer wg.Done()
// 		for i := 0; i < 5; i++ {
// 			ob.PlaceLimitOrder(orderbook.Buy, clientID, decimal.NewFromInt(10), decimal.NewFromInt(10))
// 		}
// 	}()
// 	wg.Add(1)
// 	go func() {
// 		fmt.Println("start2")
// 		defer wg.Done()
// 		for i := 0; i < 51; i++ {
// 			ob.PlaceMarketOrder(orderbook.Sell, clientID, decimal.NewFromInt(1))
// 		}
// 	}()
// 	wg.Wait()
// 	fmt.Printf("max: %v\n", max)
// 	// fmt.Printf("order: %v\n", ob.activeOrders[orderID].Value.logs)
// 	fmt.Printf("orderside volume %v\n", ob.BidVolume())
// }

// // func TestMarketOrderPartialFill(t *testing.T) {
// // 	ob := NewService("AAPL")
// // 	clientID := ulid.Make()
// // 	_, err := ob.PlaceLimitOrder(orderbook.Buy, clientID, decimal.NewFromInt(10), decimal.NewFromInt(10))
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	ob.PlaceMarketOrder(orderbook.Sell, clientID, decimal.NewFromInt(15))
// // 	// Log(fmt.Sprintf("order side volume %v\n", ob.asks.volume))
// // 	// Log(fmt.Sprintf("order queue volume %v\n", oq.volume))
// // 	assert(t, ob.AskVolume(), "5")
// // 	assert(t, ob.bids.volume.String(), "0")
// // 	assert(t, ob.asks.depth, 0)
// // }

// // func TestMarketOrderVolumeAndDepth(t *testing.T) {
// // 	ob := NewService("AAPL")
// // 	clientID := ulid.Make()
// // 	_, err := ob.PlaceMarketOrder(orderbook.Sell, clientID, decimal.NewFromInt(15))
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	_, err = ob.PlaceMarketOrder(orderbook.Sell, clientID, decimal.NewFromInt(5))
// // 	if err != nil {
// // 		t.Error(err)
// // 	}

// // 	// Log(fmt.Sprintf("order side volume %v\n", ob.asks.volume))
// // 	// Log(fmt.Sprintf("order queue volume %v\n", oq.volume))
// // 	assert(t, ob.asks.volume.String(), "20")
// // 	assert(t, ob.bids.volume.String(), "0")
// // 	assert(t, ob.asks.depth, 0)
// // 	_, err = ob.PlaceLimitOrder(orderbook.Sell, clientID, decimal.NewFromInt(10), decimal.NewFromInt(10))
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	// Log(fmt.Sprintf("order side volume %v\n", ob.asks.volume))
// // 	assert(t, ob.bids.volume.String(), "0")
// // 	assert(t, ob.asks.volume.String(), "30")
// // 	assert(t, ob.asks.depth, 1)
// // 	_, err = ob.PlaceLimitOrder(orderbook.Sell, clientID, decimal.NewFromInt(10), decimal.NewFromInt(10))
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	// Log(fmt.Sprintf("order side volume %v\n", ob.asks.volume))
// // 	assert(t, ob.bids.volume.String(), "0")
// // 	assert(t, ob.asks.volume.String(), "40")
// // 	assert(t, ob.asks.depth, 1)
// // 	_, err = ob.PlaceLimitOrder(orderbook.Sell, clientID, decimal.NewFromInt(15), decimal.NewFromInt(9))
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	// Log(fmt.Sprintf("order side volume %v\n", ob.asks.volume))
// // 	assert(t, ob.bids.volume.String(), "0")
// // 	assert(t, ob.asks.volume.String(), "55")
// // 	assert(t, ob.asks.depth, 2)
// // 	// assert(t, ob.asks.priceTree.Len(), 2)
// // 	assert(t, ob.bids.depth, 0)
// // }

// // func TestLimitOrderFilling(t *testing.T) {
// // 	ob := NewService("AAPL")
// // 	clientID := ulid.Make()
// // 	_, err := ob.PlaceMarketOrder(orderbook.Sell, clientID, decimal.NewFromInt(15))
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	assert(t, ob.asks.volume.String(), "15")
// // 	assert(t, ob.bids.volume.String(), "0")
// // 	assert(t, ob.asks.depth, 0)
// // 	_, err = ob.PlaceLimitOrder(orderbook.Buy, clientID, decimal.NewFromInt(5), decimal.NewFromInt(10))
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	assert(t, ob.asks.volume.String(), "10")
// // 	assert(t, ob.bids.volume.String(), "0")
// // 	assert(t, ob.asks.depth, 0)

// // 	_, err = ob.PlaceLimitOrder(orderbook.Buy, clientID, decimal.NewFromInt(5), decimal.NewFromInt(20))
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	assert(t, ob.asks.volume.String(), "5")
// // 	assert(t, ob.bids.volume.String(), "0")
// // 	assert(t, ob.asks.depth, 0)

// // 	_, err = ob.PlaceLimitOrder(orderbook.Buy, clientID, decimal.NewFromInt(6), decimal.NewFromInt(30))
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	assert(t, ob.asks.volume.String(), "0")
// // 	assert(t, ob.bids.volume.String(), "1")
// // 	assert(t, ob.asks.depth, 0)
// // 	assert(t, ob.bids.depth, 1)

// // 	_, err = ob.PlaceLimitOrder(orderbook.Sell, clientID, decimal.NewFromInt(2), decimal.NewFromInt(50))
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	assert(t, ob.asks.volume.String(), "2")
// // 	assert(t, ob.bids.volume.String(), "1")
// // 	assert(t, ob.asks.depth, 1)
// // 	assert(t, ob.bids.depth, 1)

// // 	_, err = ob.PlaceLimitOrder(orderbook.Buy, clientID, decimal.NewFromInt(10), decimal.NewFromInt(51))
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	assert(t, ob.asks.volume.String(), "0")
// // 	assert(t, ob.bids.volume.String(), "9")
// // 	assert(t, ob.asks.depth, 0)
// // 	assert(t, ob.bids.depth, 2)
// // }

// func assert(t *testing.T, a, b any) {
// 	if !reflect.DeepEqual(a, b) {
// 		t.Errorf("%+v != %+v", a, b)
// 	}
// }

func TestSimulateStockMarketFluctuations(t *testing.T) {
	fmt.Println("start test")

	cfg, err := config.Load("../.env")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	rdb := redis.NewRedis(cfg.Rdb)

	db, err := db.New(cfg.DB)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	obRepo := orderbook.NewRepository(db.DB)
	ob := orderbook.NewService("AAPL", obRepo, rdb)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			time.Sleep(time.Millisecond * 10)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		clientID := ulid.Make()
		for i := 0; i < 4000; i++ {
			price := decimal.NewFromInt(rand.Int63n(10) + 10)
			quantity := decimal.NewFromInt(rand.Int63n(10) + 1)
			ob.PlaceLimitOrder(orderbook.Buy, clientID, quantity, price)
			// time.Sleep(time.Millisecond * 10)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		clientID := ulid.Make()
		for i := 0; i < 5000; i++ {
			quantity := decimal.NewFromInt(rand.Int63n(3) + 1)
			ob.PlaceMarketOrder(orderbook.Buy, clientID, quantity)
			// time.Sleep(time.Millisecond * 10)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		clientID := ulid.Make()
		for i := 0; i < 5000; i++ {
			price := decimal.NewFromInt(rand.Int63n(10) + 10)
			quantity := decimal.NewFromInt(rand.Int63n(10) + 1)
			ob.PlaceLimitOrder(orderbook.Sell, clientID, quantity, price)
			// time.Sleep(time.Millisecond * 10)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		clientID := ulid.Make()
		for i := 0; i < 5500; i++ {
			quantity := decimal.NewFromInt(rand.Int63n(3) + 1)
			ob.PlaceMarketOrder(orderbook.Sell, clientID, quantity)
			// time.Sleep(time.Millisecond * 10)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		clientID := ulid.Make()
		for i := 0; i < 5000; i++ {
			price := decimal.NewFromInt(rand.Int63n(10) + 10)
			quantity := decimal.NewFromInt(rand.Int63n(10) + 1)
			ob.PlaceLimitOrder(orderbook.Sell, clientID, quantity, price)
			// time.Sleep(time.Millisecond * 10)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		clientID := ulid.Make()
		for i := 0; i < 5000; i++ {
			quantity := decimal.NewFromInt(rand.Int63n(3) + 1)
			ob.PlaceMarketOrder(orderbook.Buy, clientID, quantity)
			// time.Sleep(time.Millisecond * 10)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		clientID := ulid.Make()
		for i := 0; i < 5000; i++ {
			price := decimal.NewFromInt(rand.Int63n(10) + 10)
			quantity := decimal.NewFromInt(rand.Int63n(10) + 1)
			ob.PlaceLimitOrder(orderbook.Sell, clientID, quantity, price)
			// time.Sleep(time.Millisecond * 10)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		clientID := ulid.Make()
		for i := 0; i < 5000; i++ {
			quantity := decimal.NewFromInt(rand.Int63n(3) + 1)
			ob.PlaceMarketOrder(orderbook.Sell, clientID, quantity)
			// time.Sleep(time.Millisecond * 10)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		clientID := ulid.Make()
		for i := 0; i < 5000; i++ {
			price := decimal.NewFromInt(rand.Int63n(10) + 10)
			quantity := decimal.NewFromInt(rand.Int63n(10) + 1)
			ob.PlaceLimitOrder(orderbook.Buy, clientID, quantity, price)
			// time.Sleep(time.Millisecond * 10)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		clientID := ulid.Make()
		for i := 0; i < 5000; i++ {
			quantity := decimal.NewFromInt(rand.Int63n(3) + 1)
			ob.PlaceMarketOrder(orderbook.Buy, clientID, quantity)
			// time.Sleep(time.Millisecond * 10)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		clientID := ulid.Make()
		for i := 0; i < 5000; i++ {
			price := decimal.NewFromInt(rand.Int63n(10) + 10)
			quantity := decimal.NewFromInt(rand.Int63n(10) + 1)
			ob.PlaceLimitOrder(orderbook.Sell, clientID, quantity, price)
			// time.Sleep(time.Millisecond * 10)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		clientID := ulid.Make()
		for i := 0; i < 5000; i++ {
			quantity := decimal.NewFromInt(rand.Int63n(3) + 1)
			ob.PlaceMarketOrder(orderbook.Sell, clientID, quantity)
			// time.Sleep(time.Millisecond * 10)
		}
	}()

	wg.Wait()
}
