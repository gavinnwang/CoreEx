package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github/wry-0313/exchange/internal/models"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	newline = []byte{'\n'}
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 5 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Symbol struct {
	streaming bool
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	user *models.User

	// A map of symbols that the client is subscribed to.
	// symbols map[string]Symbol

	// A map of subscriptions that the client has. Each value is a cancel channel to close the subscription.
	subscriptions map[string]chan bool

	// Websocket dependencies.
	ws *WebSocket

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.closeSubscriptions()
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("Failed to set read deadline: %v", err)
	}
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Printf("Failed to set read deadline to new time: %v", err)
		}
		return nil
	})
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		handleMessage(c, msg)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("Failed to set write deadline: %v", err)
			}
			if !ok {
				// The hub closed the channel.
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Printf("Failed to write message: %v", err)
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(message); err != nil {
				log.Printf("Failed to set write message: %v", err)
			}

			// Add queued messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				if _, err := w.Write(newline); err != nil {
					log.Printf("Failed to set write message: %v", err)
				}
				if _, err := w.Write(<-c.send); err != nil {
					log.Printf("Failed to set write message: %v", err)
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("Failed to set write deadline: %v", err)
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte("hi")); err != nil {
				return
			}
		}
	}
}


// subscribe subscribes a client to a Redis symbol channel
func (c *Client) subscribe(symbol string) {
	rdb := c.ws.rdb
	pubsub := rdb.Subscribe(context.Background(), symbol)
	defer pubsub.Close()

	cancel := make(chan bool)
	c.subscriptions[symbol] = cancel

	ch := pubsub.Channel()
	fmt.Printf("Channel created for symbol %v\n", symbol)
	for {
		select {
		case msg := <-ch:
			log.Printf("Received message from channel %v: %v", symbol, msg.Payload)
			// Forward messages received from pubsub channel to client
			c.send <- []byte(msg.Payload)
		case <-cancel:
			fmt.Printf("Cancelling subscription %v\n", symbol)
			return
		}
	}
}

func (c *Client) closeSubscriptions() {
	for _, cancel := range c.subscriptions {
		cancel <- true
	}
}

func handleMessage(c *Client, msg []byte) {
	// Identify message event
	var msgReq Request
	err := json.Unmarshal(msg, &msgReq)
	if err != nil {
		closeConnection(c, websocket.CloseInvalidFramePayloadData, CloseReasonBadEvent)
		return
	}

	switch msgReq.Event {
	case "":
		closeConnection(c, websocket.CloseInvalidFramePayloadData, CloseReasonBadEvent)
		return
	case EventStreamSymbolInfo:
		handleStreamSymbolInfo(c, msgReq)
	default:
		closeConnection(c, websocket.CloseInvalidFramePayloadData, CloseReasonUnsupportedEvent)
		return
	}

}

// closeConnection is a helper function to write a control message with a status code and close reason text.
func closeConnection(c *Client, statusCode int, text string) {
	err := c.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(statusCode, text), time.Now().Add(writeWait))
	if err != nil {
		log.Printf("Failed to write control message: %v", err)
	}
}
