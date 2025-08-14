package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub manages WebSocket connections with simple event handling
type Hub struct {
	// Connection management
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client

	// Simple room management: room_name -> clients
	chatRooms  map[string]map[*Client]bool
	roomsMutex sync.RWMutex

	// Post subscribers: post_id -> clients
	postSubscribers map[string]map[*Client]bool
	postMutex       sync.RWMutex

	// WebSocket upgrader
	upgrader websocket.Upgrader
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:         make(map[*Client]bool),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		chatRooms:       make(map[string]map[*Client]bool),
		postSubscribers: make(map[string]map[*Client]bool),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// Run starts the hub's main event loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.handleClientRegister(client)

		case client := <-h.unregister:
			h.handleClientUnregister(client)
		}
	}
}

// handleClientRegister adds a new client
func (h *Hub) handleClientRegister(client *Client) {
	h.clients[client] = true
	log.Printf("✅ Client %s connected", client.id)
}

// handleClientUnregister removes a client from all rooms and subscriptions
func (h *Hub) handleClientUnregister(client *Client) {
	if _, ok := h.clients[client]; ok {
		// Remove from clients
		delete(h.clients, client)
		close(client.send)

		// Remove from all chat rooms
		h.roomsMutex.Lock()
		for roomName, roomClients := range h.chatRooms {
			if _, exists := roomClients[client]; exists {
				delete(roomClients, client)
				if len(roomClients) == 0 {
					delete(h.chatRooms, roomName)
				}
			}
		}
		h.roomsMutex.Unlock()

		// Remove from all post subscriptions
		h.postMutex.Lock()
		for postID, postClients := range h.postSubscribers {
			if _, exists := postClients[client]; exists {
				delete(postClients, client)
				if len(postClients) == 0 {
					delete(h.postSubscribers, postID)
				}
			}
		}
		h.postMutex.Unlock()

		log.Printf("❌ Client %s disconnected", client.id)
	}
}
