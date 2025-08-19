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

	// Unified channel system: channel_name -> clients
	channels      map[string]map[*Client]bool
	channelsMutex sync.RWMutex

	// User connection mapping: username -> client (for direct messaging)
	userConnections map[string]*Client
	usersMutex      sync.RWMutex

	// Legacy support (for backward compatibility)
	chatRooms       map[string]map[*Client]bool
	postSubscribers map[string]map[*Client]bool
	roomsMutex      sync.RWMutex
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
		channels:        make(map[string]map[*Client]bool),
		userConnections: make(map[string]*Client),
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
	
	// Initialize client subscriptions map
	client.subsMutex.Lock()
	if client.subscriptions == nil {
		client.subscriptions = make(map[string]bool)
	}
	client.subsMutex.Unlock()
	
	// Register user connection for direct messaging
	if client.username != "" {
		h.usersMutex.Lock()
		// Remove old connection if exists
		if oldClient, exists := h.userConnections[client.username]; exists {
			oldClient.mutex.Lock()
			oldClient.isConnected = false
			oldClient.mutex.Unlock()
		}
		h.userConnections[client.username] = client
		h.usersMutex.Unlock()
		log.Printf("👤 User %s registered for direct messaging", client.username)
	}
	
	log.Printf("✅ Client %s connected", client.id)
}

// handleClientUnregister removes a client from all rooms and subscriptions
func (h *Hub) handleClientUnregister(client *Client) {
	if _, ok := h.clients[client]; ok {
		// Remove from clients
		delete(h.clients, client)
		close(client.send)

		// Remove from unified channel system
		h.channelsMutex.Lock()
		for channelName, channelClients := range h.channels {
			if _, exists := channelClients[client]; exists {
				delete(channelClients, client)
				if len(channelClients) == 0 {
					delete(h.channels, channelName)
				}
			}
		}
		h.channelsMutex.Unlock()

		// Remove from user connections
		if client.username != "" {
			h.usersMutex.Lock()
			if h.userConnections[client.username] == client {
				delete(h.userConnections, client.username)
			}
			h.usersMutex.Unlock()
		}

		// Legacy cleanup - Remove from all chat rooms
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

		// Legacy cleanup - Remove from all post subscriptions
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
