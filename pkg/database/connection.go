package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func NewDatabase() (*DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./chat.db"
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &DB{DB: db}
	
	if err := database.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Insert sample data for demo
	if err := database.InsertSampleData(); err != nil {
		log.Printf("Warning: Failed to insert sample data: %v", err)
	}

	log.Println("Database connected successfully")
	return database, nil
}

func (db *DB) migrate() error {
	// Messages table
	createMessagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		content TEXT NOT NULL,
		room_id TEXT NOT NULL,
		type TEXT NOT NULL DEFAULT 'message',
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Posts table
	createPostsTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		author_id TEXT NOT NULL,
		author_name TEXT NOT NULL,
		comment_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Comments table
	createCommentsTable := `
	CREATE TABLE IF NOT EXISTS comments (
		id TEXT PRIMARY KEY,
		post_id TEXT NOT NULL,
		content TEXT NOT NULL,
		author_id TEXT NOT NULL,
		author_name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`

	// Create indexes
	createIndexes := `
	CREATE INDEX IF NOT EXISTS idx_messages_room_id ON messages(room_id);
	CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages(timestamp);
	CREATE INDEX IF NOT EXISTS idx_messages_room_timestamp ON messages(room_id, timestamp);
	
	CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);
	CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts(author_id);
	
	CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
	CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at);
	CREATE INDEX IF NOT EXISTS idx_comments_author_id ON comments(author_id);`

	// Execute migrations
	tables := []string{createMessagesTable, createPostsTable, createCommentsTable}
	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	if _, err := db.Exec(createIndexes); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}