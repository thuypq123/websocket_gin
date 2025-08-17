package handlers

import (
	"websocket/internal/repository"
	"websocket/internal/websocket"

	"github.com/gin-gonic/gin"
)

// WebSocketHandler handles WebSocket connections for all domains (chat, comments, posts)
type WebSocketHandler struct {
	hub         *websocket.Hub
	messageRepo *repository.MessageRepository
	commentRepo *repository.CommentRepository
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(
	hub *websocket.Hub,
	messageRepo *repository.MessageRepository,
	commentRepo *repository.CommentRepository,
) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		messageRepo: messageRepo,
		commentRepo: commentRepo,
	}
}

// HandleWebSocket handles WebSocket connections for all domains
// Supports: chat messages, post comments, room management
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Delegate to Hub for connection management
	h.hub.HandleWebSocket(c)
}

// GetStats returns WebSocket connection statistics
func (h *WebSocketHandler) GetStats(c *gin.Context) {
	stats := h.hub.GetStats()
	c.JSON(200, gin.H{
		"websocket_stats": stats,
		"supported_events": []string{
			"JOIN_ROOM",
			"CHAT_MESSAGE",
			"POST_COMMENT",
		},
		"domains": []string{
			"chat",
			"comments",
			"rooms",
		},
	})
}
