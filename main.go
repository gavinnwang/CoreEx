package main

import (
	"fmt"
	"time"

	"github.com/wry0313/crypto-exchange/client"
	"github.com/wry0313/crypto-exchange/server"
)

func main() {
	go server.StartServer()

	time.Sleep(1 * time.Second)

	c := client.NewClient()

	bidParams := &client.PlaceLimitOrderParams{
		UserID: 1,
		Bid:    true,
		Price:  10000,
		Size:   1000000000000000000,
	}

	go func() {
		for {
			resp, err := c.PlaceLimitOrder(bidParams)
			if err != nil {
				panic(err)
			}

			fmt.Printf("order id => %d\n", resp.OrderID)
			time.Sleep(1 * time.Second)
		}
	}()

	askParams := &client.PlaceLimitOrderParams{
		UserID: 1,
		Bid:    false,
		Price:  8000,
		Size:   1000000000000000000,
	}
	go func() {
		for {
			resp, err := c.PlaceLimitOrder(askParams)
			if err != nil {
				panic(err)
			}
			fmt.Printf("order id => %d\n", resp.OrderID)
			time.Sleep(1 * time.Second)
		}
	}()	
	select {}
}