package comments

import (
	"fmt"
	"regexp"
)

// Validator handles validation for comment events
type Validator struct {
	postIDRegex *regexp.Regexp
}

// NewValidator creates a new comments validator
func NewValidator() *Validator {
	return &Validator{
		postIDRegex: regexp.MustCompile(`^[a-zA-Z0-9_-]+$`),
	}
}

// ValidatePostComment validates a post comment event
func (v *Validator) ValidatePostComment(event *PostCommentEvent) error {
	if event.PostID == "" {
		return fmt.Errorf("post_id is required for comment")
	}
	if event.Comment == "" {
		return fmt.Errorf("comment content is required")
	}
	if len(event.Comment) > 2000 {
		return fmt.Errorf("comment too long (max 2000 characters)")
	}
	if len(event.PostID) > 100 {
		return fmt.Errorf("post_id too long (max 100 characters)")
	}
	if !v.postIDRegex.MatchString(event.PostID) {
		return fmt.Errorf("invalid post_id format (only alphanumeric, dash, underscore allowed)")
	}
	return nil
}
