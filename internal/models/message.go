package models

import "time"

type Message struct {
	ID        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Content   string    `json:"content" db:"content"`
	RoomID    string    `json:"room_id" db:"room_id"`
	Type      string    `json:"type" db:"type"` // "message", "join", "leave"
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}


type Client struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

