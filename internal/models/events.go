package models

import "time"

// Base event structure for WebSocket communication
type WSEvent struct {
	Type      string      `json:"type"`       // "chat", "comment", "post"  
	Action    string      `json:"action"`     // "send", "create", "update", "delete"
	Data      interface{} `json:"data"`       // Actual payload
	UserID    string      `json:"user_id"`    // User performing action
	Username  string      `json:"username"`   // Username for display
	RoomID    string      `json:"room_id"`    // Chat room or post ID
	Timestamp time.Time   `json:"timestamp"`  // Event timestamp
	EventID   string      `json:"event_id"`   // Unique event identifier
}

// Chat event payload
type ChatEventData struct {
	Message Message `json:"message"`
}

// Comment event payload  
type CommentEventData struct {
	Comment Comment `json:"comment"`
	PostID  string  `json:"post_id"`
}

// Post event payload
type PostEventData struct {
	Post Post `json:"post"`
}

// Event types constants
const (
	// Event Types
	EventTypeChat    = "chat"
	EventTypeComment = "comment" 
	EventTypePost    = "post"
	
	// Actions
	ActionSend   = "send"     // Chat message
	ActionCreate = "create"   // Create comment/post
	ActionUpdate = "update"   // Update comment/post  
	ActionDelete = "delete"   // Delete comment/post
	ActionJoin   = "join"     // Join chat room
	ActionLeave  = "leave"    // Leave chat room
)

// Room types for organizing connections
type RoomType string

const (
	RoomTypeChat RoomType = "chat"
	RoomTypePost RoomType = "post"
)

// Enhanced room structure
type Room struct {
	ID       string             `json:"id"`
	Name     string             `json:"name"`
	Type     RoomType           `json:"type"`     // chat or post
	Clients  map[string]*Client `json:"clients"`
	Metadata map[string]string  `json:"metadata"` // Extra info like post title, etc
}