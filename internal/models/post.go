package models

import "time"

type Post struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	AuthorID    string    `json:"author_id" db:"author_id"`
	AuthorName  string    `json:"author_name" db:"author_name"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	CommentCount int      `json:"comment_count" db:"comment_count"`
}

type Comment struct {
	ID        string    `json:"id" db:"id"`
	PostID    string    `json:"post_id" db:"post_id"`
	Content   string    `json:"content" db:"content"`
	AuthorID  string    `json:"author_id" db:"author_id"`
	AuthorName string   `json:"author_name" db:"author_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}