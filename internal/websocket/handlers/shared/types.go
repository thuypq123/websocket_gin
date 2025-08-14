package shared

// ClientInterface defines what handlers need from a client
type ClientInterface interface {
	GetUsername() string
	GetID() string
	GetHub() HubInterface
	SendError(message string)
}

// HubInterface defines what handlers need from the hub
type HubInterface interface {
	JoinChatRoom(client ClientInterface, roomName string)
	SubscribeToPost(client ClientInterface, postID string)
	BroadcastToChatRoom(roomName string, event interface{})
	BroadcastToPostSubscribers(postID string, event interface{})
	SendToClient(client ClientInterface, event interface{}) error
}

// Event interface - all events must implement this
type Event interface {
	GetType() string
	GetUser() string
}
