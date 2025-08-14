package repository

import (
	"fmt"
	"time"

	"websocket/internal/models"
	"websocket/pkg/database"
)

type PostRepository struct {
	db *database.DB
}

func NewPostRepository(db *database.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) CreatePost(post *models.Post) error {
	query := `
		INSERT INTO posts (id, title, content, author_id, author_name, comment_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	if post.CreatedAt.IsZero() {
		post.CreatedAt = now
	}
	if post.UpdatedAt.IsZero() {
		post.UpdatedAt = now
	}

	_, err := r.db.Exec(query,
		post.ID,
		post.Title,
		post.Content,
		post.AuthorID,
		post.AuthorName,
		post.CommentCount,
		post.CreatedAt.Format("2006-01-02 15:04:05"),
		post.UpdatedAt.Format("2006-01-02 15:04:05"),
	)

	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

func (r *PostRepository) GetPostByID(id string) (*models.Post, error) {
	query := `
		SELECT id, title, content, author_id, author_name, comment_count, created_at, updated_at
		FROM posts 
		WHERE id = ?
	`

	post := &models.Post{}
	var createdAtStr, updatedAtStr string

	err := r.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.AuthorName,
		&post.CommentCount,
		&createdAtStr,
		&updatedAtStr,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	// Parse timestamps
	if post.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr); err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}
	if post.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr); err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}

	return post, nil
}

func (r *PostRepository) GetAllPosts(limit, offset int) ([]*models.Post, error) {
	query := `
		SELECT id, title, content, author_id, author_name, comment_count, created_at, updated_at
		FROM posts 
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		var createdAtStr, updatedAtStr string

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.AuthorName,
			&post.CommentCount,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}

		// Parse timestamps
		if post.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr); err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}
		if post.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr); err != nil {
			return nil, fmt.Errorf("failed to parse updated_at: %w", err)
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate posts: %w", err)
	}

	return posts, nil
}

func (r *PostRepository) UpdatePost(post *models.Post) error {
	query := `
		UPDATE posts 
		SET title = ?, content = ?, updated_at = ?
		WHERE id = ?
	`

	post.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		post.Title,
		post.Content,
		post.UpdatedAt.Format("2006-01-02 15:04:05"),
		post.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	return nil
}

func (r *PostRepository) DeletePost(id string) error {
	query := `DELETE FROM posts WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

func (r *PostRepository) IncrementCommentCount(postID string) error {
	query := `UPDATE posts SET comment_count = comment_count + 1 WHERE id = ?`

	_, err := r.db.Exec(query, postID)
	if err != nil {
		return fmt.Errorf("failed to increment comment count: %w", err)
	}

	return nil
}

func (r *PostRepository) DecrementCommentCount(postID string) error {
	query := `UPDATE posts SET comment_count = comment_count - 1 WHERE id = ? AND comment_count > 0`

	_, err := r.db.Exec(query, postID)
	if err != nil {
		return fmt.Errorf("failed to decrement comment count: %w", err)
	}

	return nil
}