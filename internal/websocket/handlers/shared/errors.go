package shared

import "fmt"

// Common error types for WebSocket handlers
var (
	ErrInvalidEventFormat = fmt.Errorf("invalid event format")
	ErrUnknownEventType   = fmt.Errorf("unknown event type")
	ErrValidationFailed   = fmt.Errorf("validation failed")
	ErrClientDisconnected = fmt.Errorf("client disconnected")
	ErrSendBufferFull     = fmt.Errorf("send buffer full")
)

// ValidationError represents a validation error with details
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// ErrorEvent represents an error response to client
type ErrorEvent struct {
	Type    string `json:"type"`           // "ERROR"
	Message string `json:"message"`        // Error message
	Code    string `json:"code,omitempty"` // Optional error code
}

// GetType returns the event type
func (e *ErrorEvent) GetType() string { return e.Type }

// GetUser returns empty string for error events
func (e *ErrorEvent) GetUser() string { return "" }

// NewErrorEvent creates a new error event
func NewErrorEvent(message string) *ErrorEvent {
	return &ErrorEvent{
		Type:    "ERROR",
		Message: message,
	}
}

// NewErrorEventWithCode creates a new error event with error code
func NewErrorEventWithCode(message, code string) *ErrorEvent {
	return &ErrorEvent{
		Type:    "ERROR",
		Message: message,
		Code:    code,
	}
}
