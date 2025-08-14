package main

import (
	"log"
	"net/http"
	"os"

	"websocket/internal/handlers"
	"websocket/internal/repository"
	"websocket/internal/websocket"
	"websocket/pkg/database"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize repositories
	messageRepo := repository.NewMessageRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	// Initialize event router with repositories
	websocket.InitializeEventRouter(messageRepo, commentRepo)

	// Initialize simple WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Setup routes
	router := handlers.SetupEnhancedRoutes(hub, messageRepo, postRepo, commentRepo)

	log.Printf("🚀 WebSocket server starting on port %s", port)
	log.Printf("📝 Visit http://localhost:%s for Posts & Comments demo", port)
	log.Printf("💬 Visit http://localhost:%s/chat for Chat demo", port)
	log.Printf("🔗 WebSocket endpoint: ws://localhost:%s/ws", port)
	log.Println("")
	log.Println("Features available:")
	log.Println("  • Real-time chat rooms")
	log.Println("  • Post creation and management")
	log.Println("  • Real-time comment system")
	log.Println("  • Event-based WebSocket architecture")
	log.Println("  • RESTful API endpoints")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
