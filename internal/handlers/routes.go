package handlers

import (
	"websocket/internal/repository"
	"websocket/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(hub *websocket.Hub, messageRepo *repository.MessageRepository) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Initialize handlers
	chatHandler := NewChatHandler(hub, messageRepo)
	websocketHandler := NewWebSocketHandler(hub, messageRepo, nil) // nil commentRepo for legacy routes

	r.GET("/", chatHandler.IndexPage)
	r.GET("/chat", chatHandler.ChatPage)
	r.GET("/ws", websocketHandler.HandleWebSocket) // Use dedicated WebSocket handler

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "WebSocket Chat API is running",
			})
		})

		// Message endpoints
		api.GET("/messages", chatHandler.GetMessages)
		api.GET("/messages/recent", chatHandler.GetRecentMessages)
	}

	return r
}
