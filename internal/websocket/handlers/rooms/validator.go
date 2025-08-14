package rooms

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator handles validation for room events
type Validator struct {
	roomNameRegex *regexp.Regexp
	reservedRooms map[string]bool
}

// NewValidator creates a new rooms validator
func NewValidator() *Validator {
	return &Validator{
		roomNameRegex: regexp.MustCompile(`^[a-zA-Z0-9_-]+$`),
		reservedRooms: map[string]bool{
			"admin":   true,
			"system":  true,
			"private": true,
		},
	}
}

// ValidateJoinRoom validates a room join event
func (v *Validator) ValidateJoinRoom(event *JoinRoomEvent) error {
	if event.Room == "" {
		return fmt.Errorf("room name is required")
	}

	roomName := strings.TrimSpace(strings.ToLower(event.Room))

	if len(roomName) < 2 {
		return fmt.Errorf("room name too short (min 2 characters)")
	}
	if len(roomName) > 30 {
		return fmt.Errorf("room name too long (max 30 characters)")
	}
	if !v.roomNameRegex.MatchString(roomName) {
		return fmt.Errorf("invalid room name format (only alphanumeric, dash, underscore allowed)")
	}
	if v.reservedRooms[roomName] {
		return fmt.Errorf("room name '%s' is reserved", roomName)
	}

	return nil
}
