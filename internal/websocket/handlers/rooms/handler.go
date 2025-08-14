package rooms

import (
	"encoding/json"
	"fmt"
	"log"

	"websocket/internal/websocket/handlers/shared"
)

// Handler handles room-related WebSocket events
type Handler struct {
	validator *Validator
}

// NewHandler creates a new rooms handler
func NewHandler() *Handler {
	return &Handler{
		validator: NewValidator(),
	}
}

// JoinRoomEvent represents a room join event
type JoinRoomEvent struct {
	Type string `json:"type"` // "JOIN_ROOM"
	Room string `json:"room"` // Room name to join
	User string `json:"user"` // Username
}

// RoomJoinedEvent represents a room joined confirmation
type RoomJoinedEvent struct {
	Type string `json:"type"` // "ROOM_JOINED"
	Room string `json:"room"` // Room that was joined
	User string `json:"user"` // Username who joined
}

// GetType returns the event type
func (e *JoinRoomEvent) GetType() string { return e.Type }

// GetUser returns the user
func (e *JoinRoomEvent) GetUser() string { return e.User }

// GetType returns the event type
func (e *RoomJoinedEvent) GetType() string { return e.Type }

// GetUser returns the user
func (e *RoomJoinedEvent) GetUser() string { return e.User }

// HandleJoinRoom processes room join requests
func (h *Handler) HandleJoinRoom(client shared.ClientInterface, messageBytes []byte) error {
	// Parse event
	var event JoinRoomEvent
	if err := json.Unmarshal(messageBytes, &event); err != nil {
		return fmt.Errorf("invalid JOIN_ROOM event: %v", err)
	}

	// Validate event
	if err := h.validator.ValidateJoinRoom(&event); err != nil {
		return err
	}

	// Set user from client if not provided
	if event.User == "" {
		event.User = client.GetUsername()
	}

	log.Printf("🏠 Client %s joining room: %s", event.User, event.Room)

	// Join the chat room
	client.GetHub().JoinChatRoom(client, event.Room)

	// Send confirmation back to client
	response := &RoomJoinedEvent{
		Type: "ROOM_JOINED",
		Room: event.Room,
		User: client.GetUsername(),
	}

	return client.GetHub().SendToClient(client, response)
}
