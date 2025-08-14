package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"websocket/internal/models"
)

// GetMockPost returns a hardcoded post for demo purposes
func GetMockPost(c *gin.Context) {
	postID := c.Param("id")
	
	// Hardcoded mock post
	mockPost := &models.Post{
		ID:          postID,
		Title:       "🚀 Demo WebSocket Real-time Comments",
		Content:     "Đây là bài post demo để test hệ thống comment real-time với WebSocket. Bạn có thể thử comment và sẽ thấy các comment xuất hiện real-time trên tất cả các kết nối WebSocket đang lắng nghe post này.\n\nTính năng:\n• Real-time comments qua WebSocket\n• Event-driven architecture\n• Room-based broadcasting\n• Persistent storage với SQLite\n\nHãy thử comment bên dưới! 💬",
		AuthorID:    "demo-user-123",
		AuthorName:  "Demo User",
		CommentCount: 3,
		CreatedAt:   time.Now().Add(-2 * time.Hour),
		UpdatedAt:   time.Now().Add(-1 * time.Hour),
	}

	c.JSON(http.StatusOK, gin.H{"post": mockPost})
}

// GetMockComments returns hardcoded comments for demo
func GetMockComments(c *gin.Context) {
	postID := c.Param("id")
	
	mockComments := []*models.Comment{
		{
			ID:         "mock-comment-1",
			PostID:     postID,
			Content:    "Bài viết rất hay! Cảm ơn bạn đã chia sẻ kiến thức về WebSocket 🎉",
			AuthorID:   "user-1",
			AuthorName: "Alice",
			CreatedAt:  time.Now().Add(-90 * time.Minute),
			UpdatedAt:  time.Now().Add(-90 * time.Minute),
		},
		{
			ID:         "mock-comment-2", 
			PostID:     postID,
			Content:    "Event-driven architecture thực sự rất powerful. Tôi đã implement tương tự và performance tăng đáng kể! 🚀",
			AuthorID:   "user-2",
			AuthorName: "Bob",
			CreatedAt:  time.Now().Add(-60 * time.Minute),
			UpdatedAt:  time.Now().Add(-60 * time.Minute),
		},
		{
			ID:         "mock-comment-3",
			PostID:     postID,
			Content:    "Có thể share source code của dự án không? Tôi muốn học thêm về WebSocket với Go 💻",
			AuthorID:   "user-3",
			AuthorName: "Charlie",
			CreatedAt:  time.Now().Add(-30 * time.Minute),
			UpdatedAt:  time.Now().Add(-30 * time.Minute),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": mockComments,
		"post_id":  postID,
		"count":    len(mockComments),
	})
}

// GetAllMockPosts returns list of mock posts
func GetAllMockPosts(c *gin.Context) {
	mockPosts := []*models.Post{
		{
			ID:          "demo-post-123",
			Title:       "🚀 Demo WebSocket Real-time Comments",
			Content:     "Đây là bài post demo để test hệ thống comment real-time...",
			AuthorID:    "demo-user-123", 
			AuthorName:  "Demo User",
			CommentCount: 3,
			CreatedAt:   time.Now().Add(-2 * time.Hour),
			UpdatedAt:   time.Now().Add(-1 * time.Hour),
		},
		{
			ID:          "demo-post-456",
			Title:       "📚 Học Go từ cơ bản đến nâng cao",
			Content:     "Series bài viết hướng dẫn học Go programming language từ A-Z...",
			AuthorID:    "author-456",
			AuthorName:  "Go Expert", 
			CommentCount: 5,
			CreatedAt:   time.Now().Add(-4 * time.Hour),
			UpdatedAt:   time.Now().Add(-3 * time.Hour),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"posts":  mockPosts,
		"limit":  10,
		"offset": 0,
		"count":  len(mockPosts),
	})
}