package server

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
		UserID    int64
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

func StartServer() {
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
	pkStr_1 := os.Getenv("USER_1_PRIVATE_KEY")
	pkStr_2 := os.Getenv("USER_2_PRIVATE_KEY")
	pkStr_3 := os.Getenv("USER_3_PRIVATE_KEY")

	user1 := NewUser(1, pkStr_1)
	ex.Users[user1.ID] = user1
	user1Balance, err := client.BalanceAt(context.Background(), crypto.PubkeyToAddress(user1.PrivateKey.PublicKey), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("user1 balance: %d\n", user1Balance)

	user2 := NewUser(2, pkStr_2)
	ex.Users[user2.ID] = user2
	user2Balance, err := client.BalanceAt(context.Background(), crypto.PubkeyToAddress(user2.PrivateKey.PublicKey), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("user2 balance: %d\n", user2Balance)

	user3 := NewUser(3, pkStr_3)
	ex.Users[user3.ID] = user3
	user3Balance, err := client.BalanceAt(context.Background(), crypto.PubkeyToAddress(user3.PrivateKey.PublicKey), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("user3 balance: %d\n", user3Balance)	

	e.GET("/book/:market", ex.handleGetBook)
	e.POST("/order", ex.handlePlaceOrder)
	e.DELETE("/order/:id", ex.cancelOrder)
	e.GET("/balance/:id", ex.handleGetUserBalance)

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

	// log.Println(balance)

	e.Start(":3000")
}

type User struct {
	ID         int64
	PrivateKey *ecdsa.PrivateKey // we need this to sign the transaction with the user's private key
}

func NewUser(id int64, privKey string) *User {
	pk, err := crypto.HexToECDSA(privKey)
	if err != nil {
		panic(err)
	}
	log.Printf("creating userID: %v\n", id)
	return &User{
		ID:         id,
		PrivateKey: pk,
	}
}

func httperrorHandler(err error, c echo.Context) {
	log.Println(err)
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

func (ex *Exchange) handleGetUserBalance(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"msg": "invalid user id"})
	}
	user, ok := ex.Users[int64(userID)]
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]any{"msg": "user not found"})
	}
	balance, err := ex.getBalance(user.PrivateKey)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"msg": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{"balance": balance})
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
				UserID:   order.UserID,
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
				UserID:   order.UserID,
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
	_, ok := ex.Users[order.UserID]
	if !ok {
		return fmt.Errorf("user not found for userID: %v", order.UserID)
	}
	ob := ex.orderbooks[market]
	ob.PlaceLimitOrder(price, order)
	log.Printf("new LIMIT order => type: [%t] | price [%.2f] | size [%.2f]\n", order.Bid, order.Limit.Price, order.Size)

	return nil
}


type PlaceOrderResponse struct {
	OrderID int64
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
			return c.JSON(http.StatusBadRequest, map[string]any{"msg": err.Error()})
		}
		resp := &PlaceOrderResponse{
			OrderID: order.ID,
		}
		return c.JSON(200, resp)
	}

	if placeOrderData.Type == MarketOrder {
		matches, matchedOrders := ex.handlePlaceMarketOrder(market, order)

		if err := ex.handleMatches(matches); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{"msg": err.Error()})
		}

		return c.JSON(200, map[string]any{"matches": matchedOrders})
	}

	return nil
}

func (ex *Exchange) handleMatches(matches []orderbook.Match) error {
	for _, match := range matches {
		fromUser, ok := ex.Users[match.Ask.UserID]
		if !ok {
			return fmt.Errorf("user not found for userID: %v", match.Ask.UserID)
		}

		toUser, ok := ex.Users[match.Bid.UserID]
		if !ok {
			return fmt.Errorf("user not found for userID: %v", match.Bid.UserID)
		}
		toAddress := crypto.PubkeyToAddress(toUser.PrivateKey.PublicKey)
		// this is only used for the fees
		// exchangePubKey := ex.PrivateKey.Public()
		// publicKeyECDSA, ok := exchangePubKey.(*ecdsa.PublicKey)
		// if !ok {
		// 	return fmt.Errorf("error casting pubic key to ECDSA")
		// }
		log.Printf("transfer from %v to %v\n", fromUser.ID, toUser.ID)
		amount := big.NewInt(int64(match.SizeFilled))
		err := transferETH(ex.Client, fromUser.PrivateKey, toAddress, amount)
		if err != nil {
			return fmt.Errorf("error transferring ETH: %v", err)
		}
	}
	return nil
}

func transferETH(client *ethclient.Client, fromPrivKey *ecdsa.PrivateKey, to common.Address, amount *big.Int) error {
	ctx := context.Background()
	publicKey := fromPrivKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting pubic key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	balance, err := client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("balance: %d\n", new(big.Int).Div(balance, big.NewInt(1000000000000000000)))
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return fmt.Errorf("error getting nonce: %v", err)
	}
	
	gasLimit := uint64(21000) // in units
	gasPrice, err :=  client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// gasPrice := big.NewInt(30000000000)
	// log.Printf("gas price: %v\n", gasPrice)

	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)
	log.Printf("amount to transfer: %v\n", amount)
	chainID := big.NewInt(1337)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromPrivKey)
	if err != nil {
		return err
	}

	balance, err = client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("balance after transferr: %d\n", new(big.Int).Div(balance, big.NewInt(1000000000000000000)))
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return fmt.Errorf("error sending transaction: %v", err)
	}
	log.Printf("transfer successful")
	return nil
}

func (ex *Exchange) getBalance(privateKey *ecdsa.PrivateKey) (*big.Int, error) {
	pubKey := privateKey.Public()
	publicKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting pubic key to ECDSA")
	}
	exchangeAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	balance, err := ex.Client.BalanceAt(context.Background(), exchangeAddress, nil)
	balanceInEth := new(big.Int).Div(balance, big.NewInt(1000000000000000000))
	if err != nil {
		return nil, fmt.Errorf("error getting balance: %v", err)
	}
	return balanceInEth, nil
}
