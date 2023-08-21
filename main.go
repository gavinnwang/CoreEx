package main

import (
	"time"

	"github.com/wry0313/crypto-exchange/client"
	"github.com/wry0313/crypto-exchange/server"
)

func main() {
	go server.StartServer()

	time.Sleep(1 * time.Second)

	c := client.NewClient()

	for {
		limitOrderParams := &client.PlaceOrderParams{
			UserID: 1,
			Bid: false,
			Price: 10_000,
			Size: 500_000_000_000_000_000,
		}

		_, err := c.PlaceLimitOrder(limitOrderParams)
		if err != nil {
			panic (err)
		}

		otherLimitOrderParams := &client.PlaceOrderParams{
			UserID: 3,
			Bid: false,
			Price: 9_000,
			Size: 500_000_000_000_000_000,
		}
		
		_, err = c.PlaceLimitOrder(otherLimitOrderParams)
		if err != nil {
			panic (err)
		}

		// fmt.Println("Limit order placed: ", resp.OrderID)

		marketOrderParams := &client.PlaceOrderParams{
			UserID: 2,
			Bid: true,
			Size: 1_000_000_000_000_000_000,
		}

		_, err = c.PlaceMarketOrder(marketOrderParams)
		if err != nil {
			panic (err)
		}

		// fmt.Println("Market order placed: ", resp.OrderID)

		time.Sleep(1 * time.Second)
	}

	select {}
}