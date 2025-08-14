package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"websocket/internal/repository"
	"websocket/internal/websocket"
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

	chatHandler := NewChatHandler(hub, messageRepo)

	r.GET("/", chatHandler.IndexPage)
	r.GET("/chat", chatHandler.ChatPage)
	r.GET("/ws", chatHandler.WebSocketEndpoint)

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"message": "WebSocket Chat API is running",
			})
		})
		
		// Message endpoints
		api.GET("/messages", chatHandler.GetMessages)
		api.GET("/messages/recent", chatHandler.GetRecentMessages)
	}

	return r
}