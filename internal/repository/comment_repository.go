package repository

import (
	"fmt"
	"time"

	"websocket/internal/models"
	"websocket/pkg/database"
)

type CommentRepository struct {
	db *database.DB
}

func NewCommentRepository(db *database.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) CreateComment(comment *models.Comment) error {
	query := `
		INSERT INTO comments (id, post_id, content, author_id, author_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	if comment.CreatedAt.IsZero() {
		comment.CreatedAt = now
	}
	if comment.UpdatedAt.IsZero() {
		comment.UpdatedAt = now
	}

	_, err := r.db.Exec(query,
		comment.ID,
		comment.PostID,
		comment.Content,
		comment.AuthorID,
		comment.AuthorName,
		comment.CreatedAt.Format("2006-01-02 15:04:05"),
		comment.UpdatedAt.Format("2006-01-02 15:04:05"),
	)

	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

func (r *CommentRepository) GetCommentsByPostID(postID string, limit, offset int) ([]*models.Comment, error) {
	query := `
		SELECT id, post_id, content, author_id, author_name, created_at, updated_at
		FROM comments 
		WHERE post_id = ?
		ORDER BY created_at ASC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		var createdAtStr, updatedAtStr string

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.Content,
			&comment.AuthorID,
			&comment.AuthorName,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}

		// Parse timestamps
		if comment.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr); err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}
		if comment.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr); err != nil {
			return nil, fmt.Errorf("failed to parse updated_at: %w", err)
		}

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate comments: %w", err)
	}

	return comments, nil
}

func (r *CommentRepository) GetRecentCommentsByPostID(postID string, limit int) ([]*models.Comment, error) {
	query := `
		SELECT id, post_id, content, author_id, author_name, created_at, updated_at
		FROM comments 
		WHERE post_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, postID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent comments: %w", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		var createdAtStr, updatedAtStr string

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.Content,
			&comment.AuthorID,
			&comment.AuthorName,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}

		// Parse timestamps
		if comment.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr); err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}
		if comment.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr); err != nil {
			return nil, fmt.Errorf("failed to parse updated_at: %w", err)
		}

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate recent comments: %w", err)
	}

	// Reverse to get chronological order (oldest first)
	for i, j := 0, len(comments)-1; i < j; i, j = i+1, j-1 {
		comments[i], comments[j] = comments[j], comments[i]
	}

	return comments, nil
}

func (r *CommentRepository) UpdateComment(comment *models.Comment) error {
	query := `
		UPDATE comments 
		SET content = ?, updated_at = ?
		WHERE id = ?
	`

	comment.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		comment.Content,
		comment.UpdatedAt.Format("2006-01-02 15:04:05"),
		comment.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

func (r *CommentRepository) DeleteComment(id string) error {
	query := `DELETE FROM comments WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

func (r *CommentRepository) GetCommentByID(id string) (*models.Comment, error) {
	query := `
		SELECT id, post_id, content, author_id, author_name, created_at, updated_at
		FROM comments 
		WHERE id = ?
	`

	comment := &models.Comment{}
	var createdAtStr, updatedAtStr string

	err := r.db.QueryRow(query, id).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.Content,
		&comment.AuthorID,
		&comment.AuthorName,
		&createdAtStr,
		&updatedAtStr,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	// Parse timestamps
	if comment.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr); err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}
	if comment.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr); err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}

	return comment, nil
}