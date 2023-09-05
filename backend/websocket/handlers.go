package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var (
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all origins
			return true
		},
	}
)

func (ws *WebSocket) HandleConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection to WebSocket
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("handler: failed to upgrade connection: %v", err)
		// logger.Errorf("handler: failed to upgrade connection: %v", err)
		// logger.Info("request:", r)
		return
	}

	client := Client{
		symbols:       make(map[string]Symbol),
		subscriptions: make(map[string]chan bool),
		conn:          conn,
		send:          make(chan []byte, 256),
		ws:            ws,
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func handleStreamSymbolInfo(c *Client, msgReq Request) {

	var params ParamsSymbol
	if err := UnmarshalParams(msgReq, &params, c); err != nil {
		return
	}
	symbol := params.Symbol

	// Create a new ticker that fires every second.
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	stopChan := make(chan struct{})

	// Run the ticker in a separate goroutine.
	go func() {
		for {
			select {
			case <-ticker.C:
				symbolInfo, err := c.ws.exchangeService.GetSymbolInfo(symbol)
				if err != nil {
					log.Printf("Failed to get symbol info: %v", err)
					sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
					return
				}

				// Broacast successful message response
				msgRes := ResponseGetSymbolInfo{
					ResponseBase: ResponseBase{
						Event:   EventStreamSymbolInfo,
						Success: true,
					},
					Result: symbolInfo,
				}

				msgResBytes, err := json.Marshal(msgRes)
				if err := handleMarshalError(err, "handleStreamSymbolInfo", c); err != nil {
					return
				}
				c.send <- msgResBytes
			case <-stopChan:
				log.Println("Stopping stream symbol info ticker")
				return
			}
		}
	}()

	// Wait for the client to close the connection or any other stop condition.
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			// Stop the ticker if the client has disconnected.
			stopChan <- struct{}{}
			return
		}
	}
}

func handleStreamMarketPrice(c *Client, msgReq Request) {

	var params ParamsSymbol
	if err := UnmarshalParams(msgReq, &params, c); err != nil {
		return
	}
	symbol := params.Symbol

	// Create a new ticker that fires every second.
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	stopChan := make(chan struct{})

	// Run the ticker in a separate goroutine.
	go func() {
		for {
			select {
			case <-ticker.C:
				// Fetch the current market price here for 'symbol'.
				// For demonstration, let's assume the price is a random value.
				price, err := c.ws.exchangeService.GetMarketPrice(symbol)
				if err != nil {
					log.Printf("Failed to get market price: %v", err)
					sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
					return
				}

				// Broacast successful message response
				msgRes := ResponseGetMarketPrice{
					ResponseBase: ResponseBase{
						Event:   EventStreamMarketPrice,
						Success: true,
					},
					Result: ResultGetMarketPrice{
						Symbol: symbol,
						Price:  price,
					},
				}

				msgResBytes, err := json.Marshal(msgRes)
				if err := handleMarshalError(err, "handleStreamMarketPrice", c); err != nil {
					return
				}
				c.send <- msgResBytes
			case <-stopChan:
				log.Println("Stopping stream price ticker")
				return
			}
		}
	}()

	// Wait for the client to close the connection or any other stop condition.
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			// Stop the ticker if the client has disconnected.
			stopChan <- struct{}{}
			return
		}
	}
}

// unmarshalParams is a helper function that unmarshals a message request's params and sends
// out a close connection message if any errors are encountered.
func UnmarshalParams(msgReq Request, v any, c *Client) error {
	err := json.Unmarshal(msgReq.Params, v)
	if err != nil {
		closeConnection(c, websocket.CloseInvalidFramePayloadData, CloseReasonBadParams)
		return err
	}
	return nil
}

// handleMarshalError checks to see if there are any errors when marshalling the WebSocket message response into JSON.
// If there is an issue, it will close the connection with an internal server error close reason.
func handleMarshalError(err error, handlerName string, c *Client) error {
	if err != nil {
		log.Printf("%s: Failed to marshal response into JSON: %v", handlerName, err)
		closeConnection(c, websocket.CloseProtocolError, CloseReasonInternalServer)
		return err
	}
	return nil
}

func (ws *WebSocket) RegisterHandlers(r chi.Router) {
	r.HandleFunc("/ws", ws.HandleConnection)
}
