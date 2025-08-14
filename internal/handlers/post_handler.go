package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"websocket/internal/models"
	"websocket/internal/repository"
)

type PostHandler struct {
	postRepo    *repository.PostRepository
	commentRepo *repository.CommentRepository
}

func NewPostHandler(postRepo *repository.PostRepository, commentRepo *repository.CommentRepository) *PostHandler {
	return &PostHandler{
		postRepo:    postRepo,
		commentRepo: commentRepo,
	}
}

// GetAllPosts retrieves all posts with pagination
func (h *PostHandler) GetAllPosts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
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

	posts, err := h.postRepo.GetAllPosts(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts":  posts,
		"limit":  limit,
		"offset": offset,
		"count":  len(posts),
	})
}

// GetPostByID retrieves a single post by ID
func (h *PostHandler) GetPostByID(c *gin.Context) {
	postID := c.Param("id")

	post, err := h.postRepo.GetPostByID(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// CreatePost creates a new post
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content" binding:"required"`
		AuthorID string `json:"author_id" binding:"required"`
		AuthorName string `json:"author_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := &models.Post{
		ID:         generateID(),
		Title:      req.Title,
		Content:    req.Content,
		AuthorID:   req.AuthorID,
		AuthorName: req.AuthorName,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := h.postRepo.CreatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"post": post})
}

// UpdatePost updates an existing post
func (h *PostHandler) UpdatePost(c *gin.Context) {
	postID := c.Param("id")
	
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.postRepo.GetPostByID(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Update fields if provided
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	if err := h.postRepo.UpdatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// DeletePost deletes a post
func (h *PostHandler) DeletePost(c *gin.Context) {
	postID := c.Param("id")

	// Check if post exists
	if _, err := h.postRepo.GetPostByID(postID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err := h.postRepo.DeletePost(postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// GetCommentsByPostID retrieves all comments for a post
func (h *PostHandler) GetCommentsByPostID(c *gin.Context) {
	postID := c.Param("id")
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

	// Check if post exists
	if _, err := h.postRepo.GetPostByID(postID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	comments, err := h.commentRepo.GetCommentsByPostID(postID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"post_id":  postID,
		"limit":    limit,
		"offset":   offset,
		"count":    len(comments),
	})
}

// GetRecentCommentsByPostID retrieves recent comments for a post
func (h *PostHandler) GetRecentCommentsByPostID(c *gin.Context) {
	postID := c.Param("id")
	limitStr := c.DefaultQuery("limit", "20")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	// Check if post exists
	if _, err := h.postRepo.GetPostByID(postID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	comments, err := h.commentRepo.GetRecentCommentsByPostID(postID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recent comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"post_id":  postID,
		"count":    len(comments),
	})
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}