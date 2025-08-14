package websocket

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket connection constants
const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

// HandleWebSocket upgrades HTTP connection to WebSocket
func (h *Hub) HandleWebSocket(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("❌ WebSocket upgrade failed: %v", err)
		return
	}

	// Create client
	username := c.Query("username")
	if username == "" {
		username = "anonymous"
	}

	client := &Client{
		hub:         h,
		conn:        conn,
		send:        make(chan []byte, 256),
		id:          generateClientID(),
		username:    username,
		isConnected: true,
	}

	// Register client and start goroutines
	h.register <- client
	go client.writePump()
	go client.readPump()
}

// readPump handles incoming WebSocket messages
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		c.mutex.Lock()
		c.isConnected = false
		c.mutex.Unlock()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle the event
		if err := c.handleEvent(messageBytes); err != nil {
			log.Printf("❌ Error handling event: %v", err)
			c.SendError(err.Error())
		}
	}
}

// writePump handles outgoing WebSocket messages
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Send sends raw bytes to the client
func (c *Client) Send(data []byte) error {
	select {
	case c.send <- data:
		return nil
	default:
		return fmt.Errorf("client send buffer is full")
	}
}
