package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a simple WebSocket client
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte

	// Client info
	id       string
	username string

	// Connection state
	isConnected bool
	mutex       sync.Mutex

	// Multi-channel subscriptions
	subscriptions map[string]bool // channel_name -> subscribed
	subsMutex     sync.RWMutex    // Thread-safe access to subscriptions
}
