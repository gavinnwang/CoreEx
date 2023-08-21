package client

import "net/http"

type Client struct {
	*http.Client
}

func NewClient() *Client {
	return &Client{
		Client: http.DefaultClient,
	}
}

func (c *Client) PlaceLimitOrder() error {
	return nil
}