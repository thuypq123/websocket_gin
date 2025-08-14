package events

import (
	"encoding/json"
	"log"
	"time"

	"websocket/internal/models"
	"websocket/internal/repository"
)

// ChatEventHandler handles chat-related WebSocket events
type ChatEventHandler struct {
	messageRepo *repository.MessageRepository
}

func NewChatEventHandler(messageRepo *repository.MessageRepository) *ChatEventHandler {
	return &ChatEventHandler{
		messageRepo: messageRepo,
	}
}

func (h *ChatEventHandler) GetEventType() string {
	return models.EventTypeChat
}

func (h *ChatEventHandler) HandleEvent(event *models.WSEvent, client Client) error {
	switch event.Action {
	case models.ActionSend:
		return h.handleSendMessage(event, client)
	case models.ActionJoin:
		return h.handleJoinRoom(event, client)
	case models.ActionLeave:
		return h.handleLeaveRoom(event, client)
	default:
		return ErrInvalidEventData{
			EventType: models.EventTypeChat,
			Reason:    "unsupported action: " + event.Action,
		}
	}
}

func (h *ChatEventHandler) handleSendMessage(event *models.WSEvent, client Client) error {
	// Parse chat event data
	var chatData models.ChatEventData
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return ErrInvalidEventData{
			EventType: models.EventTypeChat,
			Reason:    "failed to marshal event data",
		}
	}

	if err := json.Unmarshal(dataBytes, &chatData); err != nil {
		return ErrInvalidEventData{
			EventType: models.EventTypeChat,
			Reason:    "invalid chat event data format",
		}
	}

	// Create message
	message := &models.Message{
		ID:        generateID(),
		Username:  client.GetUsername(),
		Content:   chatData.Message.Content,
		RoomID:    client.GetRoomID(),
		Type:      "message",
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	// Save to database
	if err := h.messageRepo.SaveMessage(message); err != nil {
		log.Printf("Failed to save chat message: %v", err)
		// Don't return error, continue with broadcast
	}

	// Create broadcast event
	broadcastEvent := &models.WSEvent{
		Type:      models.EventTypeChat,
		Action:    models.ActionSend,
		Data:      models.ChatEventData{Message: *message},
		UserID:    client.GetUserID(),
		Username:  client.GetUsername(),
		RoomID:    client.GetRoomID(),
		Timestamp: time.Now(),
		EventID:   generateID(),
	}

	// Broadcast to chat room
	return client.BroadcastToRoom(models.RoomTypeChat, client.GetRoomID(), broadcastEvent)
}

func (h *ChatEventHandler) handleJoinRoom(event *models.WSEvent, client Client) error {
	// Create join message
	message := &models.Message{
		ID:        generateID(),
		Username:  client.GetUsername(),
		Content:   client.GetUsername() + " joined the chat",
		RoomID:    client.GetRoomID(),
		Type:      "join",
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	// Save to database
	if err := h.messageRepo.SaveMessage(message); err != nil {
		log.Printf("Failed to save join message: %v", err)
	}

	// Create broadcast event
	broadcastEvent := &models.WSEvent{
		Type:      models.EventTypeChat,
		Action:    models.ActionJoin,
		Data:      models.ChatEventData{Message: *message},
		UserID:    client.GetUserID(),
		Username:  client.GetUsername(),
		RoomID:    client.GetRoomID(),
		Timestamp: time.Now(),
		EventID:   generateID(),
	}

	// Broadcast to chat room
	return client.BroadcastToRoom(models.RoomTypeChat, client.GetRoomID(), broadcastEvent)
}

func (h *ChatEventHandler) handleLeaveRoom(event *models.WSEvent, client Client) error {
	// Create leave message
	message := &models.Message{
		ID:        generateID(),
		Username:  client.GetUsername(),
		Content:   client.GetUsername() + " left the chat",
		RoomID:    client.GetRoomID(),
		Type:      "leave",
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	// Save to database
	if err := h.messageRepo.SaveMessage(message); err != nil {
		log.Printf("Failed to save leave message: %v", err)
	}

	// Create broadcast event
	broadcastEvent := &models.WSEvent{
		Type:      models.EventTypeChat,
		Action:    models.ActionLeave,
		Data:      models.ChatEventData{Message: *message},
		UserID:    client.GetUserID(),
		Username:  client.GetUsername(),
		RoomID:    client.GetRoomID(),
		Timestamp: time.Now(),
		EventID:   generateID(),
	}

	// Broadcast to chat room
	return client.BroadcastToRoom(models.RoomTypeChat, client.GetRoomID(), broadcastEvent)
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}