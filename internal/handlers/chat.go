package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"websocket/internal/repository"
	"websocket/internal/websocket"
)

type ChatHandler struct {
	hub         *websocket.Hub
	messageRepo *repository.MessageRepository
}

func NewChatHandler(hub *websocket.Hub, messageRepo *repository.MessageRepository) *ChatHandler {
	return &ChatHandler{
		hub:         hub,
		messageRepo: messageRepo,
	}
}

func (h *ChatHandler) WebSocketEndpoint(c *gin.Context) {
	h.hub.HandleWebSocket(c)
}

func (h *ChatHandler) ChatPage(c *gin.Context) {
	c.HTML(http.StatusOK, "chat.html", gin.H{
		"title": "WebSocket Chat",
	})
}

func (h *ChatHandler) IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Welcome to Chat",
	})
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	roomID := c.Query("room")
	if roomID == "" {
		roomID = "general"
	}

	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	messages, err := h.messageRepo.GetMessagesByRoom(roomID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"room_id":  roomID,
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *ChatHandler) GetRecentMessages(c *gin.Context) {
	roomID := c.Query("room")
	if roomID == "" {
		roomID = "general"
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	messages, err := h.messageRepo.GetRecentMessagesByRoom(roomID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recent messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"room_id":  roomID,
		"count":    len(messages),
	})
}