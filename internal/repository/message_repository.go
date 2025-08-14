package repository

import (
	"fmt"
	"time"

	"websocket/internal/models"
	"websocket/pkg/database"
)

type MessageRepository struct {
	db *database.DB
}

func NewMessageRepository(db *database.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (r *MessageRepository) SaveMessage(message *models.Message) error {
	query := `
		INSERT INTO messages (id, username, content, room_id, type, timestamp, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	if message.CreatedAt.IsZero() {
		message.CreatedAt = now
	}
	if message.Timestamp.IsZero() {
		message.Timestamp = now
	}

	_, err := r.db.Exec(query,
		message.ID,
		message.Username,
		message.Content,
		message.RoomID,
		message.Type,
		message.Timestamp.Format("2006-01-02 15:04:05"),
		message.CreatedAt.Format("2006-01-02 15:04:05"),
	)

	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	return nil
}

func (r *MessageRepository) GetMessagesByRoom(roomID string, limit int, offset int) ([]*models.Message, error) {
	query := `
		SELECT id, username, content, room_id, type, timestamp, created_at
		FROM messages 
		WHERE room_id = ? 
		ORDER BY timestamp ASC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, roomID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		var timestampStr, createdAtStr string

		err := rows.Scan(
			&message.ID,
			&message.Username,
			&message.Content,
			&message.RoomID,
			&message.Type,
			&timestampStr,
			&createdAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		// Parse timestamps - try multiple formats
		message.Timestamp, err = parseFlexibleTimestamp(timestampStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}
		message.CreatedAt, err = parseFlexibleTimestamp(createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}

		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate messages: %w", err)
	}

	return messages, nil
}

func (r *MessageRepository) GetRecentMessagesByRoom(roomID string, limit int) ([]*models.Message, error) {
	query := `
		SELECT id, username, content, room_id, type, timestamp, created_at
		FROM messages 
		WHERE room_id = ? 
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, roomID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		var timestampStr, createdAtStr string

		err := rows.Scan(
			&message.ID,
			&message.Username,
			&message.Content,
			&message.RoomID,
			&message.Type,
			&timestampStr,
			&createdAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		// Parse timestamps - try multiple formats
		message.Timestamp, err = parseFlexibleTimestamp(timestampStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}
		message.CreatedAt, err = parseFlexibleTimestamp(createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}

		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate recent messages: %w", err)
	}

	// Reverse the slice to get chronological order (oldest first)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *MessageRepository) DeleteOldMessages(roomID string, olderThan time.Time) error {
	query := `
		DELETE FROM messages 
		WHERE room_id = ? AND timestamp < ?
	`

	_, err := r.db.Exec(query, roomID, olderThan.Format("2006-01-02 15:04:05"))
	if err != nil {
		return fmt.Errorf("failed to delete old messages: %w", err)
	}

	return nil
}

func (r *MessageRepository) GetMessageCount(roomID string) (int, error) {
	query := `SELECT COUNT(*) FROM messages WHERE room_id = ?`

	var count int
	err := r.db.QueryRow(query, roomID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get message count: %w", err)
	}

	return count, nil
}

// parseFlexibleTimestamp tries multiple timestamp formats
func parseFlexibleTimestamp(timestampStr string) (time.Time, error) {
	// Common timestamp formats
	formats := []string{
		"2006-01-02 15:04:05",      // Custom format
		time.RFC3339,               // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05Z",     // ISO 8601 UTC
		"2006-01-02T15:04:05.000Z", // ISO 8601 with milliseconds
		time.DateTime,              // 2006-01-02 15:04:05
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timestampStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse timestamp '%s' with any known format", timestampStr)
}
