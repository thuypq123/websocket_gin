package events

import "fmt"

// Custom error types for event handling
type ErrHandlerNotFound struct {
	EventType string
}

func (e ErrHandlerNotFound) Error() string {
	return fmt.Sprintf("no handler found for event type: %s", e.EventType)
}

type ErrInvalidEventData struct {
	EventType string
	Reason    string
}

func (e ErrInvalidEventData) Error() string {
	return fmt.Sprintf("invalid event data for %s: %s", e.EventType, e.Reason)
}

type ErrUnauthorized struct {
	Action string
	UserID string
}

func (e ErrUnauthorized) Error() string {
	return fmt.Sprintf("user %s not authorized to perform action: %s", e.UserID, e.Action)
}

type ErrResourceNotFound struct {
	ResourceType string
	ResourceID   string
}

func (e ErrResourceNotFound) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.ResourceType, e.ResourceID)
}