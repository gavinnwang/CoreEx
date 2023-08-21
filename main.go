package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"

	"github.com/wry0313/crypto-exchange/orderbook"
)

const (
	MarketOrder OrderType = "MARKET"
	LimitOrder  OrderType = "LIMIT"

	MarketETH Market = "ETH"
)

type (
	OrderType string
	Market    string

	PlaceOrderRequest struct {
		UserID int64
		Type   OrderType // Limit or Market
		Bid    bool
		Size   float64
		Price  float64
		Market Market
	}

	// only limit order because market order doesn't have a price
	Order struct {
		ID        int64
		Price     float64
		Size      float64
		Bid       bool
		Timestamp int64
	}

	OrderbookData struct {
		TotalBidVolume float64
		TotalAskVolume float64
		Asks           []*Order
		Bids           []*Order
	}

	Exchange struct {
		Client     *ethclient.Client
		Users      map[int64]*User
		orders     map[int64]int64 // orderID -> userID
		PrivateKey *ecdsa.PrivateKey
		orderbooks map[Market]*orderbook.Orderbook
	}

	CancelOrderRequest struct {
		Bid bool
		ID  int64
	}

	MatchedOrder struct {
		ID    int64
		Price float64
		Size  float64
	}
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()
	e.HTTPErrorHandler = httperrorHandler

	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}

	ex, err := NewExchange(os.Getenv("EXCHANGE_PRIVATE_KEY"), client)
	if err != nil {
		log.Fatal(err)
	}
	pkStr := "61880c63e6f3088fce0a401f68fc2343625107223b13529bb13cc567c2de2f38"

	// pk, err := crypto.HexToECDSA(pkStr)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	user := NewUser(8, pkStr)
	ex.Users[user.ID] = user

	e.GET("/book/:market", ex.handleGetBook)
	e.POST("/order", ex.handlePlaceOrder)
	e.DELETE("/order/:id", ex.cancelOrder)

	// // address := common.HexToAddress("0xE19B27EcF741284AE7e3fF5F8aba026266ba25F6")

	// privateKey, err := crypto.HexToECDSA(os.Getenv("EXCHANGE_PRIVATE_KEY"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// publicKey := privateKey.Public()
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	// if !ok {
	// 	log.Fatal("cfnnot assert type: publicKey is not of type *ecdsa.PublicKey")
	// }

	// fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// value := big.NewInt(123456) // in wei (1 eth)

	// gasLimit := uint64(21000) // in units
	// gasPrice, err := client.SuggestGasPrice(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")

	// tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// chainID, err := client.NetworkID(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err = client.SendTransaction(context.Background(), signedTx); err != nil {
	// 	log.Fatal(err)
	// }

	// balance, err := client.BalanceAt(context.Background(), toAddress, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(balance)

	e.Start(":3000")
}

type User struct {
	ID         int64
	PrivateKey *ecdsa.PrivateKey // we need this to sign the transaction with the user's private key
}

func NewUser(ID int64, privKey string) *User {
	pk, err := crypto.HexToECDSA(privKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("creating userID: %v\n", ID)
	return &User{
		ID:         ID,
		PrivateKey: pk,
	}
}

func httperrorHandler(err error, c echo.Context) {
	fmt.Println(err)
}

func NewExchange(privateKey string, client *ethclient.Client) (*Exchange, error) {
	orderbooks := make(map[Market]*orderbook.Orderbook)
	orderbooks[MarketETH] = orderbook.NewOrderbook()

	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	return &Exchange{
		Client:     client,
		Users:      make(map[int64]*User),
		orders:     make(map[int64]int64),
		PrivateKey: pk,
		orderbooks: orderbooks,
	}, nil
}

func (ex *Exchange) handleGetBook(c echo.Context) error {
	market := Market(c.Param("market"))
	ob, ok := ex.orderbooks[market]
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]any{"msg": "market not found"})
	}

	orderbookData := OrderbookData{
		TotalBidVolume: ob.BidTotalVolume(),
		TotalAskVolume: ob.AskTotalVolume(),
		Asks:           []*Order{},
		Bids:           []*Order{},
	}

	for _, limit := range ob.Asks() {
		for _, order := range limit.Orders {
			o := Order{
				ID:        order.ID,
				Price:     limit.Price,
				Size:      order.Size,
				Bid:       order.Bid,
				Timestamp: order.Timestamp,
			}
			orderbookData.Asks = append(orderbookData.Asks, &o)
		}
	}

	for _, limit := range ob.Bids() {
		for _, order := range limit.Orders {
			o := Order{
				ID:        order.ID,
				Price:     limit.Price,
				Size:      order.Size,
				Bid:       order.Bid,
				Timestamp: order.Timestamp,
			}
			orderbookData.Bids = append(orderbookData.Bids, &o)
		}
	}

	return c.JSON(http.StatusOK, orderbookData)
}

func (ex *Exchange) cancelOrder(c echo.Context) error {
	idStr := c.Param("id") // fetching param is always string
	id, _ := strconv.Atoi(idStr)

	ob := ex.orderbooks[MarketETH]
	order := ob.Orders[int64(id)]
	ob.CancelOrder(order)

	return c.JSON(200, map[string]any{"msg": "order deleted"})
}

func (ex *Exchange) handlePlaceMarketOrder(
	market Market,
	order *orderbook.Order,
) ([]orderbook.Match, []*MatchedOrder) {
	ob := ex.orderbooks[market]
	matches := ob.PlaceMarketOrder(order)
	matchedOrders := make([]*MatchedOrder, len(matches))

	isBid := false
	if order.Bid {
		isBid = true
	}

	for i := 0; i < len(matchedOrders); i++ {
		id := matches[i].Bid.ID
		if isBid {
			id = matches[i].Ask.ID
		}
		matchedOrders[i] = &MatchedOrder{
			ID:    id,
			Size:  matches[i].SizeFilled,
			Price: matches[i].Price,
		}
	}

	return matches, matchedOrders
}

func (ex *Exchange) handlePlaceLimitOrder(
	market Market,
	price float64,
	order *orderbook.Order,
) error {
	user, ok := ex.Users[order.UserID]
	if !ok {
		return fmt.Errorf("user not found for userID: %v", order.UserID)
	}
	exchangePubKey := ex.PrivateKey.Public()
	publicKeyECDSA, ok := exchangePubKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting pubic key to ECDSA")
	}
	toAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	amount := big.NewInt(int64(order.Size))

	fmt.Println(toAddress.String())
	fmt.Printf("amount to transfer: %v\n", amount)
	fmt.Printf("user: %v\n", user)
	fmt.Printf("user private key: %v\n", user.PrivateKey)
	fmt.Printf("exchange private key: %v\n", ex.Client)

	// transfer from user to exchange if cannot transfer send back order
	err := transferETH(ex.Client, user.PrivateKey, toAddress, amount)
	if err != nil {
		return err
	}
	fmt.Println("transfer successful add order to orderbook")
	ob := ex.orderbooks[market]
	ob.PlaceLimitOrder(price, order)

	reserve, err := ex.getReserve()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("reserve in exchange after transfering limit order: %v\n", reserve)
	return nil
}

func (ex *Exchange) handlePlaceOrder(c echo.Context) error {
	var placeOrderData PlaceOrderRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return err
	}

	market := Market(placeOrderData.Market)
	order := orderbook.NewOrder(placeOrderData.Bid, placeOrderData.Size, placeOrderData.UserID)

	if placeOrderData.Type == LimitOrder {
		if err := ex.handlePlaceLimitOrder(market, placeOrderData.Price, order); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{"error msg": err.Error()})
		}
		return c.JSON(200, map[string]any{"msg": "limit order placed"})
	}

	if placeOrderData.Type == MarketOrder {
		matches, matchedOrders := ex.handlePlaceMarketOrder(market, order)

		if err := ex.handleMatches(matches); err != nil {
			return err
		}

		return c.JSON(200, map[string]any{"matches": matchedOrders})
	}

	return nil
}

func (ex *Exchange) handleMatches(matches []orderbook.Match) error {
	return nil
}
