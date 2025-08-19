package subscription

import (
	"fmt"
	"strings"
)

// Validator handles validation for subscription events
type Validator struct{}

// NewValidator creates a new subscription validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateSubscribe validates a subscription event
func (v *Validator) ValidateSubscribe(event *SubscribeEvent) error {
	if len(event.Channels) == 0 {
		return fmt.Errorf("at least one channel is required for subscription")
	}

	if len(event.Channels) > 50 {
		return fmt.Errorf("too many channels requested (max 50)")
	}

	// Validate each channel name
	for _, channelName := range event.Channels {
		if err := v.validateChannelName(channelName); err != nil {
			return fmt.Errorf("invalid channel '%s': %v", channelName, err)
		}
	}

	return nil
}

// ValidateUnsubscribe validates an unsubscription event
func (v *Validator) ValidateUnsubscribe(event *UnsubscribeEvent) error {
	if len(event.Channels) == 0 {
		return fmt.Errorf("at least one channel is required for unsubscription")
	}

	if len(event.Channels) > 50 {
		return fmt.Errorf("too many channels requested (max 50)")
	}

	// Validate each channel name
	for _, channelName := range event.Channels {
		if err := v.validateChannelName(channelName); err != nil {
			return fmt.Errorf("invalid channel '%s': %v", channelName, err)
		}
	}

	return nil
}

// validateChannelName validates a single channel name
func (v *Validator) validateChannelName(channelName string) error {
	if channelName == "" {
		return fmt.Errorf("channel name cannot be empty")
	}

	if len(channelName) > 100 {
		return fmt.Errorf("channel name too long (max 100 characters)")
	}

	// Check format: type:identifier
	parts := strings.SplitN(channelName, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("channel name must be in format 'type:identifier'")
	}

	channelType := parts[0]
	identifier := parts[1]

	// Validate channel type
	validTypes := map[string]bool{
		"chat":    true,
		"post":    true,
		"user":    true,
		"system":  true,
		"private": true,
	}

	if !validTypes[channelType] {
		return fmt.Errorf("unknown channel type: %s", channelType)
	}

	// Validate identifier
	if identifier == "" {
		return fmt.Errorf("channel identifier cannot be empty")
	}

	if len(identifier) > 50 {
		return fmt.Errorf("channel identifier too long (max 50 characters)")
	}

	// Additional validation based on channel type
	switch channelType {
	case "chat":
		return v.validateChatChannel(identifier)
	case "post":
		return v.validatePostChannel(identifier)
	case "user":
		return v.validateUserChannel(identifier)
	case "system":
		return v.validateSystemChannel(identifier)
	case "private":
		return v.validatePrivateChannel(identifier)
	}

	return nil
}

// validateChatChannel validates chat channel identifiers
func (v *Validator) validateChatChannel(identifier string) error {
	// Chat rooms should be alphanumeric with hyphens and underscores
	if !isValidIdentifier(identifier) {
		return fmt.Errorf("chat room name can only contain letters, numbers, hyphens, and underscores")
	}
	return nil
}

// validatePostChannel validates post channel identifiers
func (v *Validator) validatePostChannel(identifier string) error {
	// Post IDs should be alphanumeric
	if !isValidIdentifier(identifier) {
		return fmt.Errorf("post ID can only contain letters, numbers, hyphens, and underscores")
	}
	return nil
}

// validateUserChannel validates user channel identifiers
func (v *Validator) validateUserChannel(identifier string) error {
	// Username validation - alphanumeric with hyphens and underscores
	if !isValidIdentifier(identifier) {
		return fmt.Errorf("username can only contain letters, numbers, hyphens, and underscores")
	}
	return nil
}

// validateSystemChannel validates system channel identifiers
func (v *Validator) validateSystemChannel(identifier string) error {
	// System channels are predefined
	validSystemChannels := map[string]bool{
		"alerts":        true,
		"announcements": true,
		"maintenance":   true,
		"updates":       true,
	}

	if !validSystemChannels[identifier] {
		return fmt.Errorf("unknown system channel: %s", identifier)
	}
	return nil
}

// validatePrivateChannel validates private channel identifiers
func (v *Validator) validatePrivateChannel(identifier string) error {
	// Private channels can be conversation IDs or user pairs
	if !isValidIdentifier(identifier) {
		return fmt.Errorf("private channel identifier can only contain letters, numbers, hyphens, and underscores")
	}
	return nil
}

// isValidIdentifier checks if an identifier contains only allowed characters
func isValidIdentifier(identifier string) bool {
	for _, char := range identifier {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return false
		}
	}
	return true
}
