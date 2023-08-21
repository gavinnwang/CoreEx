package main

import (

	"github.com/wry0313/crypto-exchange/server"
)

func main() {
	go server.StartServer()

	// time.Sleep(1 * time.Second)

	// c := client.NewClient()

	// bidParams := &client.PlaceLimitOrderParams{
	// 	UserID: 8,
	// 	Bid:    true,
	// 	Price:  10000,
	// 	Size:   1000000000000000000,
	// }

	// go func() {
	// 	for {
	// 		if err := c.PlaceLimitOrder(bidParams); err != nil {
	// 			panic(err)
	// 		}
	// 		time.Sleep(3 * time.Second)
	// 	}
	// }()

	// askParams := &client.PlaceLimitOrderParams{
	// 	UserID: 8,
	// 	Bid:    false,
	// 	Price:  8000,
	// 	Size:   1000000000000000000,
	// }
	// go func() {
	// 	for {
	// 		if err := c.PlaceLimitOrder(askParams); err != nil {
	// 			panic(err)
	// 		}
	// 		time.Sleep(3 * time.Second)
	// 	}
	// }()	
	select {}
}
