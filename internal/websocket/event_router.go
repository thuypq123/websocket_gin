package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"websocket/internal/repository"
	"websocket/internal/websocket/handlers/chat"
	"websocket/internal/websocket/handlers/comments"
	"websocket/internal/websocket/handlers/rooms"
	"websocket/internal/websocket/handlers/shared"
	"websocket/internal/websocket/handlers/subscription"
)

// eventRouter will be initialized with repositories
var eventRouter *EventRouter

// EventRouter handles routing of WebSocket events to appropriate handlers
type EventRouter struct {
	chatHandler         *chat.Handler
	commentHandler      *comments.Handler
	roomHandler         *rooms.Handler
	subscriptionHandler *subscription.Handler
}

// InitializeEventRouter initializes the global event router with repositories
func InitializeEventRouter(messageRepo *repository.MessageRepository, commentRepo *repository.CommentRepository) {
	eventRouter = &EventRouter{
		chatHandler:         chat.NewHandler(messageRepo),
		commentHandler:      comments.NewHandler(commentRepo),
		roomHandler:         rooms.NewHandler(),
		subscriptionHandler: subscription.NewHandler(),
	}
}

// routeEvent routes incoming events to appropriate handlers
func (r *EventRouter) routeEvent(client shared.ClientInterface, messageBytes []byte) error {
	// Parse to get event type
	var baseEvent struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(messageBytes, &baseEvent); err != nil {
		return fmt.Errorf("invalid event format: %v", err)
	}

	log.Printf("📨 Routing event type: %s from client %s", baseEvent.Type, client.GetUsername())

	// Route to appropriate handler using constants
	switch baseEvent.Type {
	// Legacy events
	case EventJoinRoom:
		return r.roomHandler.HandleJoinRoom(client, messageBytes)
	case EventChatMessage:
		return r.chatHandler.HandleChatMessage(client, messageBytes)
	case EventPostComment:
		return r.commentHandler.HandlePostComment(client, messageBytes)

	// New subscription events
	case EventSubscribe:
		return r.subscriptionHandler.HandleSubscribe(client, messageBytes)
	case EventUnsubscribe:
		return r.subscriptionHandler.HandleUnsubscribe(client, messageBytes)
	case EventListSubscriptions:
		return r.subscriptionHandler.HandleListSubscriptions(client, messageBytes)

	default:
		return fmt.Errorf("unknown event type: %s", baseEvent.Type)
	}
}

// HandleEventForTesting exposes event handling for testing purposes
func HandleEventForTesting(client shared.ClientInterface, messageBytes []byte) error {
	if eventRouter == nil {
		return fmt.Errorf("event router not initialized")
	}
	return eventRouter.routeEvent(client, messageBytes)
}
