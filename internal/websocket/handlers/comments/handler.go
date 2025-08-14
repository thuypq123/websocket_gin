package comments

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"websocket/internal/models"
	"websocket/internal/repository"
	"websocket/internal/websocket/handlers/shared"
)

// Handler handles comment-related WebSocket events
type Handler struct {
	validator         *Validator
	commentRepository *repository.CommentRepository
}

// NewHandler creates a new comments handler
func NewHandler(commentRepo *repository.CommentRepository) *Handler {
	return &Handler{
		validator:         NewValidator(),
		commentRepository: commentRepo,
	}
}

// PostCommentEvent represents a post comment event
type PostCommentEvent struct {
	Type    string `json:"type"`    // "POST_COMMENT"
	PostID  string `json:"post_id"` // Target post ID
	User    string `json:"user"`    // Commenter username
	Comment string `json:"comment"` // Comment content
}

// GetType returns the event type
func (e *PostCommentEvent) GetType() string { return e.Type }

// GetUser returns the user
func (e *PostCommentEvent) GetUser() string { return e.User }

// HandlePostComment processes post comment events with database persistence
func (h *Handler) HandlePostComment(client shared.ClientInterface, messageBytes []byte) error {
	// Parse event
	var event PostCommentEvent
	if err := json.Unmarshal(messageBytes, &event); err != nil {
		return fmt.Errorf("invalid POST_COMMENT event: %v", err)
	}

	// Validate event
	if err := h.validator.ValidatePostComment(&event); err != nil {
		return err
	}

	// Set user from client if not provided
	if event.User == "" {
		event.User = client.GetUsername()
	}

	log.Printf("📝 Processing comment from %s on post %s: %s", event.User, event.PostID, event.Comment)

	// STEP 1: Save comment to database first
	comment := &models.Comment{
		ID:         generateCommentID(),
		PostID:     event.PostID,
		AuthorName: event.User,
		Content:    event.Comment,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := h.commentRepository.CreateComment(comment); err != nil {
		log.Printf("❌ Failed to save comment to database: %v", err)
		return fmt.Errorf("failed to save comment: %v", err)
	}

	log.Printf("💾 Comment saved to database with ID: %s", comment.ID)

	// STEP 2: Subscribe to this post if not already subscribed
	client.GetHub().SubscribeToPost(client, event.PostID)

	// STEP 3: Only broadcast after successful DB save
	client.GetHub().BroadcastToPostSubscribers(event.PostID, &event)

	log.Printf("📡 Comment broadcasted to post %s subscribers", event.PostID)
	return nil
}

// generateCommentID creates a unique comment ID
func generateCommentID() string {
	return fmt.Sprintf("comment_%d", time.Now().UnixNano())
}
