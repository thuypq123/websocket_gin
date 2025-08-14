package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"websocket/internal/models"
	"websocket/internal/repository"
	"websocket/internal/websocket/handlers/shared"
)

// Handler handles chat-related WebSocket events
type Handler struct {
	validator         *Validator
	messageRepository *repository.MessageRepository
}

// NewHandler creates a new chat handler
func NewHandler(messageRepo *repository.MessageRepository) *Handler {
	return &Handler{
		validator:         NewValidator(),
		messageRepository: messageRepo,
	}
}

// ChatMessageEvent represents a chat message event
type ChatMessageEvent struct {
	Type    string `json:"type"`    // "CHAT_MESSAGE"
	Room    string `json:"room"`    // Target room
	User    string `json:"user"`    // Sender username
	Message string `json:"message"` // Message content
}

// GetType returns the event type
func (e *ChatMessageEvent) GetType() string { return e.Type }

// GetUser returns the user
func (e *ChatMessageEvent) GetUser() string { return e.User }

// HandleChatMessage processes chat message events with database persistence
func (h *Handler) HandleChatMessage(client shared.ClientInterface, messageBytes []byte) error {
	// Parse event
	var event ChatMessageEvent
	if err := json.Unmarshal(messageBytes, &event); err != nil {
		return fmt.Errorf("invalid CHAT_MESSAGE event: %v", err)
	}

	// Validate event
	if err := h.validator.ValidateChatMessage(&event); err != nil {
		return err
	}

	// Set user from client if not provided
	if event.User == "" {
		event.User = client.GetUsername()
	}

	log.Printf("💬 Processing chat message from %s in room %s: %s", event.User, event.Room, event.Message)

	// STEP 1: Save to database first
	now := time.Now()
	message := &models.Message{
		ID:        generateMessageID(),
		Username:  event.User,
		Content:   event.Message,
		RoomID:    event.Room,
		Type:      "chat",
		Timestamp: now,
		CreatedAt: now,
	}

	if err := h.messageRepository.SaveMessage(message); err != nil {
		log.Printf("❌ Failed to save message to database: %v", err)
		return fmt.Errorf("failed to save message: %v", err)
	}

	log.Printf("💾 Message saved to database with ID: %s", message.ID)

	// STEP 2: Only broadcast after successful DB save
	client.GetHub().BroadcastToChatRoom(event.Room, &event)

	log.Printf("📡 Message broadcasted to room %s", event.Room)
	return nil
}

// generateMessageID creates a unique message ID
func generateMessageID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}
