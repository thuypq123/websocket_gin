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
	// Legacy methods (for backward compatibility)
	JoinChatRoom(client ClientInterface, roomName string)
	SubscribeToPost(client ClientInterface, postID string)
	BroadcastToChatRoom(roomName string, event interface{})
	BroadcastToPostSubscribers(postID string, event interface{})
	SendToClient(client ClientInterface, event interface{}) error
	
	// New channel-based methods
	SubscribeToChannel(client ClientInterface, channelName string) error
	UnsubscribeFromChannel(client ClientInterface, channelName string) error
	BroadcastToChannel(channelName string, event interface{}) error
	GetClientSubscriptions(client ClientInterface) []string
	SendToUser(username string, event interface{}) error
}

// Event interface - all events must implement this
type Event interface {
	GetType() string
	GetUser() string
}
