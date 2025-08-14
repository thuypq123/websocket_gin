package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"websocket/internal/repository"
	"websocket/internal/websocket"
	"websocket/internal/websocket/handlers/comments"
	"websocket/internal/websocket/handlers/shared"

	"github.com/gin-gonic/gin"
)

type SimpleChatHandler struct {
	hub         *websocket.Hub
	messageRepo *repository.MessageRepository
}

func NewSimpleChatHandler(hub *websocket.Hub, messageRepo *repository.MessageRepository) *SimpleChatHandler {
	return &SimpleChatHandler{
		hub:         hub,
		messageRepo: messageRepo,
	}
}

// IndexPage serves the main index page
func (h *SimpleChatHandler) IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Simple WebSocket Chat",
	})
}

// ChatPage serves the chat page
func (h *SimpleChatHandler) ChatPage(c *gin.Context) {
	room := c.DefaultQuery("room", "general")
	username := c.Query("username")

	c.HTML(http.StatusOK, "chat.html", gin.H{
		"title":    "Chat Room",
		"room":     room,
		"username": username,
	})
}

// HandleWebSocket handles WebSocket connections
func (h *SimpleChatHandler) HandleWebSocket(c *gin.Context) {
	h.hub.HandleWebSocket(c)
}

// GetRecentMessages returns recent messages for a room from database
func (h *SimpleChatHandler) GetRecentMessages(c *gin.Context) {
	room := c.Param("room")
	if room == "" {
		room = c.DefaultQuery("room", "general")
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100 // Cap at 100 messages
	}

	// Fetch messages from database
	messages, err := h.messageRepo.GetRecentMessagesByRoom(room, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch messages from database",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"room":     room,
		"count":    len(messages),
		"note":     "Messages fetched from database",
	})
}

// GetStats returns simple hub statistics
func (h *SimpleChatHandler) GetStats(c *gin.Context) {
	stats := h.hub.GetStats()
	c.JSON(http.StatusOK, stats)
}

// SendTestMessage sends a test message via proper WebSocket flow (with DB persistence)
func (h *SimpleChatHandler) SendTestMessage(c *gin.Context) {
	room := c.DefaultQuery("room", "general")
	message := c.DefaultQuery("message", "Test message from server")
	user := c.DefaultQuery("user", "system")

	// Create a mock client for testing
	mockClient := &MockClient{
		username: user,
		hub:      h.hub,
	}

	// Create WebSocket event JSON
	eventJSON := fmt.Sprintf(`{
		"type": "CHAT_MESSAGE",
		"room": "%s",
		"user": "%s", 
		"message": "%s"
	}`, room, user, message)

	// Process through proper WebSocket handler flow (includes DB save)
	if err := mockClient.handleEvent([]byte(eventJSON)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Test message sent via WebSocket flow (saved to DB)",
		"event": gin.H{
			"type":    "CHAT_MESSAGE",
			"room":    room,
			"user":    user,
			"message": message,
		},
	})
}

// SendTestComment sends a test comment (for debugging)
func (h *SimpleChatHandler) SendTestComment(c *gin.Context) {
	postID := c.DefaultQuery("post_id", "1")
	comment := c.DefaultQuery("comment", "Test comment from server")
	user := c.DefaultQuery("user", "system")

	event := &comments.PostCommentEvent{
		Type:    websocket.EventPostComment,
		PostID:  postID,
		User:    user,
		Comment: comment,
	}

	h.hub.BroadcastToPostSubscribers(postID, event)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Test comment sent",
		"event":   event,
	})
}

// MockClient implements ClientInterface for testing
type MockClient struct {
	username string
	hub      *websocket.Hub
}

func (m *MockClient) GetUsername() string         { return m.username }
func (m *MockClient) GetID() string               { return "mock_client_test" }
func (m *MockClient) GetHub() shared.HubInterface { return m.hub }
func (m *MockClient) SendError(message string) {
	log.Printf("Mock client error: %s", message)
}

func (m *MockClient) handleEvent(messageBytes []byte) error {
	// This will call the websocket event handling which includes DB persistence
	return websocket.HandleEventForTesting(m, messageBytes)
}
