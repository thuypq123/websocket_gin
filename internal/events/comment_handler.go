package events

import (
	"encoding/json"
	"log"
	"time"

	"websocket/internal/models"
	"websocket/internal/repository"
)

// CommentEventHandler handles comment-related WebSocket events
type CommentEventHandler struct {
	commentRepo *repository.CommentRepository
	postRepo    *repository.PostRepository
}

func NewCommentEventHandler(commentRepo *repository.CommentRepository, postRepo *repository.PostRepository) *CommentEventHandler {
	return &CommentEventHandler{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

func (h *CommentEventHandler) GetEventType() string {
	return models.EventTypeComment
}

func (h *CommentEventHandler) HandleEvent(event *models.WSEvent, client Client) error {
	switch event.Action {
	case models.ActionCreate:
		return h.handleCreateComment(event, client)
	case models.ActionUpdate:
		return h.handleUpdateComment(event, client)
	case models.ActionDelete:
		return h.handleDeleteComment(event, client)
	default:
		return ErrInvalidEventData{
			EventType: models.EventTypeComment,
			Reason:    "unsupported action: " + event.Action,
		}
	}
}

func (h *CommentEventHandler) handleCreateComment(event *models.WSEvent, client Client) error {
	// Parse comment event data
	var commentData models.CommentEventData
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return ErrInvalidEventData{
			EventType: models.EventTypeComment,
			Reason:    "failed to marshal event data",
		}
	}

	if err := json.Unmarshal(dataBytes, &commentData); err != nil {
		return ErrInvalidEventData{
			EventType: models.EventTypeComment,
			Reason:    "invalid comment event data format",
		}
	}

	// Validate post exists
	if _, err := h.postRepo.GetPostByID(commentData.PostID); err != nil {
		return ErrResourceNotFound{
			ResourceType: "Post",
			ResourceID:   commentData.PostID,
		}
	}

	// Create comment
	comment := &models.Comment{
		ID:         generateID(),
		PostID:     commentData.PostID,
		Content:    commentData.Comment.Content,
		AuthorID:   client.GetUserID(),
		AuthorName: client.GetUsername(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save to database
	if err := h.commentRepo.CreateComment(comment); err != nil {
		log.Printf("Failed to create comment: %v", err)
		return err
	}

	// Increment post comment count
	if err := h.postRepo.IncrementCommentCount(commentData.PostID); err != nil {
		log.Printf("Failed to increment comment count: %v", err)
	}

	// Create broadcast event
	broadcastEvent := &models.WSEvent{
		Type:      models.EventTypeComment,
		Action:    models.ActionCreate,
		Data:      models.CommentEventData{Comment: *comment, PostID: commentData.PostID},
		UserID:    client.GetUserID(),
		Username:  client.GetUsername(),
		RoomID:    commentData.PostID, // Use post ID as room ID
		Timestamp: time.Now(),
		EventID:   generateID(),
	}

	// Broadcast to post room (everyone viewing this post)
	return client.BroadcastToRoom(models.RoomTypePost, commentData.PostID, broadcastEvent)
}

func (h *CommentEventHandler) handleUpdateComment(event *models.WSEvent, client Client) error {
	// Parse comment event data
	var commentData models.CommentEventData
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return ErrInvalidEventData{
			EventType: models.EventTypeComment,
			Reason:    "failed to marshal event data",
		}
	}

	if err := json.Unmarshal(dataBytes, &commentData); err != nil {
		return ErrInvalidEventData{
			EventType: models.EventTypeComment,
			Reason:    "invalid comment event data format",
		}
	}

	// Get existing comment
	existingComment, err := h.commentRepo.GetCommentByID(commentData.Comment.ID)
	if err != nil {
		return ErrResourceNotFound{
			ResourceType: "Comment",
			ResourceID:   commentData.Comment.ID,
		}
	}

	// Check authorization (only author can update)
	if existingComment.AuthorID != client.GetUserID() {
		return ErrUnauthorized{
			Action: "update comment",
			UserID: client.GetUserID(),
		}
	}

	// Update comment
	existingComment.Content = commentData.Comment.Content
	if err := h.commentRepo.UpdateComment(existingComment); err != nil {
		log.Printf("Failed to update comment: %v", err)
		return err
	}

	// Create broadcast event
	broadcastEvent := &models.WSEvent{
		Type:      models.EventTypeComment,
		Action:    models.ActionUpdate,
		Data:      models.CommentEventData{Comment: *existingComment, PostID: existingComment.PostID},
		UserID:    client.GetUserID(),
		Username:  client.GetUsername(),
		RoomID:    existingComment.PostID,
		Timestamp: time.Now(),
		EventID:   generateID(),
	}

	// Broadcast to post room
	return client.BroadcastToRoom(models.RoomTypePost, existingComment.PostID, broadcastEvent)
}

func (h *CommentEventHandler) handleDeleteComment(event *models.WSEvent, client Client) error {
	// Parse comment event data
	var commentData models.CommentEventData
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return ErrInvalidEventData{
			EventType: models.EventTypeComment,
			Reason:    "failed to marshal event data",
		}
	}

	if err := json.Unmarshal(dataBytes, &commentData); err != nil {
		return ErrInvalidEventData{
			EventType: models.EventTypeComment,
			Reason:    "invalid comment event data format",
		}
	}

	// Get existing comment
	existingComment, err := h.commentRepo.GetCommentByID(commentData.Comment.ID)
	if err != nil {
		return ErrResourceNotFound{
			ResourceType: "Comment",
			ResourceID:   commentData.Comment.ID,
		}
	}

	// Check authorization (only author can delete)
	if existingComment.AuthorID != client.GetUserID() {
		return ErrUnauthorized{
			Action: "delete comment",
			UserID: client.GetUserID(),
		}
	}

	// Delete comment
	if err := h.commentRepo.DeleteComment(commentData.Comment.ID); err != nil {
		log.Printf("Failed to delete comment: %v", err)
		return err
	}

	// Decrement post comment count
	if err := h.postRepo.DecrementCommentCount(existingComment.PostID); err != nil {
		log.Printf("Failed to decrement comment count: %v", err)
	}

	// Create broadcast event
	broadcastEvent := &models.WSEvent{
		Type:      models.EventTypeComment,
		Action:    models.ActionDelete,
		Data:      models.CommentEventData{Comment: *existingComment, PostID: existingComment.PostID},
		UserID:    client.GetUserID(),
		Username:  client.GetUsername(),
		RoomID:    existingComment.PostID,
		Timestamp: time.Now(),
		EventID:   generateID(),
	}

	// Broadcast to post room
	return client.BroadcastToRoom(models.RoomTypePost, existingComment.PostID, broadcastEvent)
}