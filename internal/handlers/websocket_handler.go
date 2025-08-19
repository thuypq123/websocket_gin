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

	// Get additional channel statistics
	allChannels := h.hub.GetAllChannels()
	channelStats := make(map[string]int)
	for _, channel := range allChannels {
		channelStats[channel] = h.hub.GetChannelSubscribers(channel)
	}

	c.JSON(200, gin.H{
		"websocket_stats": stats,
		"channel_stats":   channelStats,
		"supported_events": []string{
			// Legacy events
			"JOIN_ROOM",
			"CHAT_MESSAGE",
			"POST_COMMENT",
			// New subscription events
			"SUBSCRIBE",
			"UNSUBSCRIBE",
			"LIST_SUBSCRIPTIONS",
		},
		"domains": []string{
			"chat",
			"comments",
			"rooms",
			"subscriptions",
		},
		"channel_types": []string{
			"chat:room_name",
			"post:post_id",
			"user:username",
			"system:alerts",
			"private:conversation_id",
		},
	})
}
