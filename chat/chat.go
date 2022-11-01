package chat

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// NewClient creates a new client
func NewClient(id int64, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:   id,
		Conn: conn,
		Send: make(chan Response, 256),
		hub:  hub,
		Room: nil,
	}
}

// Client goroutine to read messages from client
func (c *Client) Read() {

	defer func() {
		c.hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var req Request
		err := c.Conn.ReadJSON(&req)
		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
		senderID := strconv.FormatInt(c.ID, 10)
		req.Sender = senderID
		c.hub.Request <- req
	}
}

// Client goroutine to write messages to client
func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			} else {
				err := c.Conn.WriteJSON(message)
				if err != nil {
					fmt.Println("Error: ", err)
					break
				}
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}

// Client closing channel to unregister client
func (c *Client) Close() {
	close(c.Send)
}

func RegisterWS(ctx *gin.Context, hub *Hub) {
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ids, _ := ctx.Get("client-counter")
	id := ids.(int)
	// Create new client connection
	client := NewClient(int64(id), ws, hub)
	reg := RegisterStruct{
		client: *client,
	}
	hub.Register <- reg
	go client.Write()
	go client.Read()
}
