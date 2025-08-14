package database

import (
	"log"
	"time"
)

func (db *DB) InsertSampleData() error {
	// Check if sample post already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM posts WHERE id = 'sample-post-123'").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Sample data already exists, skipping insertion")
		return nil
	}

	// Insert sample post
	now := time.Now().Format("2006-01-02 15:04:05")
	
	insertPost := `
		INSERT INTO posts (id, title, content, author_id, author_name, comment_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = db.Exec(insertPost,
		"sample-post-123",
		"🚀 Hướng dẫn WebSocket với Go và Gin",
		"Bài viết này sẽ hướng dẫn chi tiết cách xây dựng hệ thống WebSocket real-time với Go và Gin framework. Chúng ta sẽ implement event-driven architecture để xử lý cả chat và comment system. Đây là một kiến trúc hiện đại, scalable và maintainable cho các ứng dụng real-time.",
		"demo-author",
		"Demo User",
		2, // comment count
		now,
		now,
	)
	if err != nil {
		return err
	}

	// Insert sample comments
	insertComment1 := `
		INSERT INTO comments (id, post_id, content, author_id, author_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	// Comment 1
	_, err = db.Exec(insertComment1,
		"sample-comment-1",
		"sample-post-123", 
		"Bài viết rất hay! Cảm ơn bạn đã chia sẻ kiến thức 🎉",
		"user-1",
		"Alice",
		now,
		now,
	)
	if err != nil {
		return err
	}

	// Comment 2  
	_, err = db.Exec(insertComment1,
		"sample-comment-2",
		"sample-post-123",
		"Event-driven architecture thực sự rất powerful. Tôi đã implement tương tự và performance tăng đáng kể! 🚀",
		"user-2", 
		"Bob",
		now,
		now,
	)
	if err != nil {
		return err
	}

	log.Println("✅ Sample data inserted successfully")
	log.Println("📝 Sample post ID: sample-post-123")
	log.Println("🔗 Visit: http://localhost:8080/post/sample-post-123")
	
	return nil
}