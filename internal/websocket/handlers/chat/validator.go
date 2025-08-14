package chat

import "fmt"

// Validator handles validation for chat events
type Validator struct{}

// NewValidator creates a new chat validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateChatMessage validates a chat message event
func (v *Validator) ValidateChatMessage(event *ChatMessageEvent) error {
	if event.Room == "" {
		return fmt.Errorf("room name is required for chat message")
	}
	if event.Message == "" {
		return fmt.Errorf("message content is required")
	}
	if len(event.Message) > 1000 {
		return fmt.Errorf("message too long (max 1000 characters)")
	}
	if len(event.Room) > 50 {
		return fmt.Errorf("room name too long (max 50 characters)")
	}
	return nil
}
