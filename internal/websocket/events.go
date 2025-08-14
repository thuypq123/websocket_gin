package websocket

// Event type constants - used by handlers
const (
	EventJoinRoom    = "JOIN_ROOM"
	EventChatMessage = "CHAT_MESSAGE"
	EventPostComment = "POST_COMMENT"
	EventRoomJoined  = "ROOM_JOINED"
	EventError       = "ERROR"
)

// Event interface - all events must implement this
type Event interface {
	GetType() string
	GetUser() string
}
