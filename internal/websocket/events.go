package websocket

// Event type constants - used by handlers
const (
	// Legacy events
	EventJoinRoom    = "JOIN_ROOM"
	EventChatMessage = "CHAT_MESSAGE"
	EventPostComment = "POST_COMMENT"
	EventRoomJoined  = "ROOM_JOINED"
	EventError       = "ERROR"

	// New subscription events
	EventSubscribe            = "SUBSCRIBE"
	EventUnsubscribe          = "UNSUBSCRIBE"
	EventListSubscriptions    = "LIST_SUBSCRIPTIONS"
	EventSubscriptionResponse = "SUBSCRIPTION_RESPONSE"
)

// Event interface - all events must implement this
type Event interface {
	GetType() string
	GetUser() string
}
