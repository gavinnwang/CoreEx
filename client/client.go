package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wry0313/crypto-exchange/server"
)

const Endpoint = "http://localhost:3000"

type PlaceLimitOrderParams struct {
	UserID int64
	Bid    bool
	Price  float64
	Size   float64
}

type Client struct {
	*http.Client
}

func NewClient() *Client {
	return &Client{
		Client: http.DefaultClient,
	}
}

func (c *Client) CancelOrder(orderID int64) error {
	e := fmt.Sprintf("%s/order/%d", Endpoint, orderID)

	req, err := http.NewRequest(http.MethodDelete, e, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error while reading response body with error status code: %v", err)
		}

		// Convert the body bytes to a string
		errorMsg := string(bodyBytes)

		// Handle or log the error message
		return fmt.Errorf("Received error response: %s", errorMsg)
	}
	return nil
}

func (c *Client) PlaceLimitOrder(p *PlaceLimitOrderParams) (*server.PlaceOrderResponse, error) {
	params := &server.PlaceOrderRequest{
		UserID: p.UserID,
		Type:   server.LimitOrder,
		Bid:    p.Bid,
		Size:   p.Size,
		Price:  p.Price,
		Market: server.MarketETH,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshalling params: %w", err)
	}

	e := Endpoint + "/order"
	req, err := http.NewRequest(http.MethodPost, e, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil,fmt.Errorf("error sending request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Error while reading response body with error status code: %v", err)
		}

		// Convert the body bytes to a string
		errorMsg := string(bodyBytes)

		// Handle or log the error message
		return nil, fmt.Errorf("Received error response: %s", errorMsg)
	}
	placeOrderResponse := &server.PlaceOrderResponse{}
	if err := json.NewDecoder(resp.Body).Decode(placeOrderResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}
	return placeOrderResponse, nil
}
