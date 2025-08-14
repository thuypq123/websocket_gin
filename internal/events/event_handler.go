package events

import (
	"websocket/internal/models"
)

// EventHandler interface for handling different types of events
type EventHandler interface {
	HandleEvent(event *models.WSEvent, client Client) error
	GetEventType() string
}

// EventManager manages all event handlers
type EventManager struct {
	handlers map[string]EventHandler
}

func NewEventManager() *EventManager {
	return &EventManager{
		handlers: make(map[string]EventHandler),
	}
}

// RegisterHandler registers a new event handler
func (em *EventManager) RegisterHandler(handler EventHandler) {
	em.handlers[handler.GetEventType()] = handler
}

// HandleEvent routes events to appropriate handlers  
func (em *EventManager) HandleEvent(event *models.WSEvent, client Client) error {
	handler, exists := em.handlers[event.Type]
	if !exists {
		return ErrHandlerNotFound{EventType: event.Type}
	}

	return handler.HandleEvent(event, client)
}

// GetHandler returns handler for specific event type
func (em *EventManager) GetHandler(eventType string) (EventHandler, bool) {
	handler, exists := em.handlers[eventType]
	return handler, exists
}

// Client represents a WebSocket client (abstracted for event handlers)
type Client interface {
	GetID() string
	GetUsername() string
	GetUserID() string
	GetRoomID() string
	Send(data []byte) error
	Broadcast(roomID string, event *models.WSEvent) error
	BroadcastToRoom(roomType models.RoomType, roomID string, event *models.WSEvent) error
}