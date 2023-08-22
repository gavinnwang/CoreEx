package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/wry0313/crypto-exchange/client"
	"github.com/wry0313/crypto-exchange/server"
)

var (
	tick   = 5 * time.Second
	myAsks = make(map[float64]int64)
	myBids = make(map[float64]int64)
)

const (
	maxOrders = 3
)

func marketOrderPlacer(c *client.Client) {
	ticker := time.NewTicker(tick)
	for {
		marketSellOrder := &client.PlaceOrderParams{
			UserID: 3, 
			Bid: false,
			Size: 1000,
		}
		sellOrderResp, err := c.PlaceMarketOrder(marketSellOrder)
		if err != nil {
			log.Println(sellOrderResp.OrderID)
		}
		marketBuyOrder := &client.PlaceOrderParams{
			UserID: 3, 
			Bid: true,
			Size: 1000,
		}
		buyOrderResp, err := c.PlaceMarketOrder(marketBuyOrder)
		if err != nil {
			log.Println(buyOrderResp.OrderID)
		}
		<-ticker.C
	}

}
func marketMakerSimple(c *client.Client) {
	ticker := time.NewTicker(tick)

	for {
		bestAsk, err := c.GetBestAsk()
		if err != nil {
			log.Println(err)
		}
		bestBid, err := c.GetBestBid()
		if err != nil {
			log.Println(err)
		}

		spread := math.Abs(bestAsk - bestBid)
		fmt.Printf("Spread: %.2f\n", spread)

		if len(myBids) < 3 {
			bidLimit := &client.PlaceOrderParams{
				UserID: 1,
				Bid:    true,
				Price:  bestBid + 100,
				Size:   1000,
			}

			bidOrderResp, err := c.PlaceLimitOrder(bidLimit)
			if err != nil {
				log.Println(bidOrderResp.OrderID)
			}

			myBids[bidLimit.Price] = bidOrderResp.OrderID
		}

		if len(myAsks) < 3 {
			askLimit := &client.PlaceOrderParams{
				UserID: 2,
				Bid:    false,
				Price:  bestAsk - 100,
				Size:   1000,
			}

			askOrderResp, err := c.PlaceLimitOrder(askLimit)
			if err != nil {
				log.Println(askOrderResp)
			}

			myAsks[askLimit.Price] = askOrderResp.OrderID
		}

		fmt.Printf("Best ask: %.2f\n", bestAsk)
		fmt.Printf("Best bid: %.2f\n", bestBid)

		<-ticker.C
	}
}

func seedMarket(c *client.Client) error {
	ask := &client.PlaceOrderParams{
		UserID: 1,
		Bid:    false,
		Price:  10_000,
		Size:   10,
	}

	bid := &client.PlaceOrderParams{
		UserID: 1,
		Bid:    true,
		Price:  9_000,
		Size:   10,
	}
	_, err := c.PlaceLimitOrder(ask)
	if err != nil {
		return err
	}
	_, err = c.PlaceLimitOrder(bid)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	go server.StartServer()

	time.Sleep(1 * time.Second)

	c := client.NewClient()

	if err := seedMarket(c); err != nil {
		panic(err)
	}

	go marketMakerSimple(c)
	time.Sleep(1 * time.Second)
	go marketOrderPlacer(c)

	// for {
	// 	limitOrderParams := &client.PlaceOrderParams{
	// 		UserID: 1,
	// 		Bid: true,
	// 		Price: 10_000,
	// 		Size: 10,
	// 	}
	// 	time.Sleep(2 * time.Second)
	// 	_, err := c.PlaceLimitOrder(limitOrderParams)
	// 	if err != nil {
	// 		panic (err)
	// 	}
	// 	time.Sleep(2 * time.Second)
	// 	otherLimitOrderParams := &client.PlaceOrderParams{
	// 		UserID: 3,
	// 		Bid: false,
	// 		Price: 12_000,
	// 		Size: 1,
	// 	}

	// 	_, err = c.PlaceLimitOrder(otherLimitOrderParams)
	// 	if err != nil {
	// 		panic (err)
	// 	}

	// 	// fmt.Println("Limit order placed: ", resp.OrderID)

	// 	marketOrderParams := &client.PlaceOrderParams{
	// 		UserID: 2,
	// 		Bid: false,
	// 		Size: 9,
	// 	}
	// 	time.Sleep(2 * time.Second)
	// 	_, err = c.PlaceMarketOrder(marketOrderParams)
	// 	if err != nil {
	// 		panic (err)
	// 	}

	// 	// fmt.Println("Market order placed: ", resp.OrderID)

	// 	time.Sleep(2 * time.Second)

	// 	bestBidPrice, err := c.GetBestBid()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	bestAskPrice, err := c.GetBestAsk()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Printf("Best bid price: %.2f\n", bestBidPrice)
	// 	fmt.Printf("Best ask price: %.2f\n", bestAskPrice)
	// }

	select {}
}
