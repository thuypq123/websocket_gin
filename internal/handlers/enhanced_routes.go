package handlers

import (
	"websocket/internal/events"
	"websocket/internal/repository"
	"websocket/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupEnhancedRoutes(
	hub *websocket.Hub,
	messageRepo *repository.MessageRepository,
	postRepo *repository.PostRepository,
	commentRepo *repository.CommentRepository,
) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// Load HTML templates and static files
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Initialize simple chat handler
	chatHandler := NewSimpleChatHandler(hub, messageRepo)
	postHandler := NewPostHandler(postRepo, commentRepo)

	// Frontend routes
	r.GET("/", chatHandler.IndexPage)
	r.GET("/chat", chatHandler.ChatPage)
	r.GET("/posts", func(c *gin.Context) {
		c.HTML(200, "posts.html", gin.H{"title": "Posts & Comments Demo"})
	})
	r.GET("/post/:id", func(c *gin.Context) {
		postID := c.Param("id")
		c.HTML(200, "post.html", gin.H{
			"title":   "Post Details",
			"post_id": postID,
		})
	})

	// WebSocket endpoint (supports both chat and post rooms)
	r.GET("/ws", chatHandler.HandleWebSocket)

	// API routes
	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Enhanced WebSocket API is running",
				"features": []string{
					"chat",
					"comments",
					"posts",
					"real-time events",
				},
			})
		})

		// Chat messages (legacy support)
		api.GET("/messages/:room", chatHandler.GetRecentMessages)
		api.GET("/messages/recent", chatHandler.GetRecentMessages)

		// Test endpoints for debugging
		api.GET("/test/message", chatHandler.SendTestMessage)
		api.GET("/test/comment", chatHandler.SendTestComment)
		api.GET("/stats", chatHandler.GetStats)

		// Posts management (using mock for demo)
		posts := api.Group("/posts")
		{
			posts.GET("", GetAllMockPosts)               // GET /api/v1/posts (MOCK)
			posts.POST("", postHandler.CreatePost)       // POST /api/v1/posts
			posts.GET("/:id", GetMockPost)               // GET /api/v1/posts/:id (MOCK)
			posts.PUT("/:id", postHandler.UpdatePost)    // PUT /api/v1/posts/:id
			posts.DELETE("/:id", postHandler.DeletePost) // DELETE /api/v1/posts/:id

			// Comments for posts (using mock for demo)
			posts.GET("/:id/comments", GetMockComments)        // GET /api/v1/posts/:id/comments (MOCK)
			posts.GET("/:id/comments/recent", GetMockComments) // GET /api/v1/posts/:id/comments/recent (MOCK)
		}
	}

	return r
}

// SetupEventManager configures all event handlers
func SetupEventManager(
	messageRepo *repository.MessageRepository,
	postRepo *repository.PostRepository,
	commentRepo *repository.CommentRepository,
) *events.EventManager {

	eventManager := events.NewEventManager()

	// Register event handlers
	chatHandler := events.NewChatEventHandler(messageRepo)
	commentHandler := events.NewCommentEventHandler(commentRepo, postRepo)

	eventManager.RegisterHandler(chatHandler)
	eventManager.RegisterHandler(commentHandler)

	return eventManager
}
