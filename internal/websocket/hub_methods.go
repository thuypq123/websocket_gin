package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"websocket/internal/websocket/handlers/shared"
)

// Make Hub implement HubInterface
func (h *Hub) JoinChatRoom(client shared.ClientInterface, roomName string) {
	// Convert interface back to concrete type for internal operations
	concreteClient, ok := client.(*Client)
	if !ok {
		log.Printf("❌ Invalid client type in JoinChatRoom")
		return
	}

	h.roomsMutex.Lock()
	defer h.roomsMutex.Unlock()

	if h.chatRooms[roomName] == nil {
		h.chatRooms[roomName] = make(map[*Client]bool)
	}
	h.chatRooms[roomName][concreteClient] = true

	log.Printf("👥 Client %s joined chat room: %s", client.GetUsername(), roomName)
}

func (h *Hub) SubscribeToPost(client shared.ClientInterface, postID string) {
	// Convert interface back to concrete type for internal operations
	concreteClient, ok := client.(*Client)
	if !ok {
		log.Printf("❌ Invalid client type in SubscribeToPost")
		return
	}

	h.postMutex.Lock()
	defer h.postMutex.Unlock()

	if h.postSubscribers[postID] == nil {
		h.postSubscribers[postID] = make(map[*Client]bool)
	}
	h.postSubscribers[postID][concreteClient] = true

	log.Printf("📝 Client %s subscribed to post: %s", client.GetUsername(), postID)
}

func (h *Hub) BroadcastToChatRoom(roomName string, event interface{}) {
	h.roomsMutex.RLock()
	roomClients := h.chatRooms[roomName]
	h.roomsMutex.RUnlock()

	if roomClients == nil {
		log.Printf("⚠️ No clients in room %s to broadcast to", roomName)
		return
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("❌ Error marshaling chat event: %v", err)
		return
	}

	for client := range roomClients {
		select {
		case client.send <- eventBytes:
		default:
			// Client's send buffer is full, remove it
			delete(roomClients, client)
			close(client.send)
		}
	}

	log.Printf("💬 Broadcasted chat message to room %s (%d clients)", roomName, len(roomClients))
}

func (h *Hub) BroadcastToPostSubscribers(postID string, event interface{}) {
	h.postMutex.RLock()
	postClients := h.postSubscribers[postID]
	h.postMutex.RUnlock()

	if postClients == nil {
		log.Printf("⚠️ No subscribers for post %s", postID)
		return
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("❌ Error marshaling comment event: %v", err)
		return
	}

	for client := range postClients {
		select {
		case client.send <- eventBytes:
		default:
			// Client's send buffer is full, remove it
			delete(postClients, client)
			close(client.send)
		}
	}

	log.Printf("📝 Broadcasted comment to post %s (%d clients)", postID, len(postClients))
}

func (h *Hub) SendToClient(client shared.ClientInterface, event interface{}) error {
	// Convert interface back to concrete type
	concreteClient, ok := client.(*Client)
	if !ok {
		return fmt.Errorf("invalid client type")
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshaling event: %v", err)
	}

	select {
	case concreteClient.send <- eventBytes:
		return nil
	default:
		return fmt.Errorf("client send buffer is full")
	}
}
