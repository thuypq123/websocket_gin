package subscription

import (
	"encoding/json"
	"fmt"
	"log"

	"websocket/internal/websocket/handlers/shared"
)

// Handler handles subscription-related WebSocket events
type Handler struct {
	validator *Validator
}

// NewHandler creates a new subscription handler
func NewHandler() *Handler {
	return &Handler{
		validator: NewValidator(),
	}
}

// SubscribeEvent represents a subscription request
type SubscribeEvent struct {
	Type     string   `json:"type"`     // "SUBSCRIBE"
	Channels []string `json:"channels"` // List of channel names to subscribe to
	User     string   `json:"user"`     // Username (auto-filled)
}

// UnsubscribeEvent represents an unsubscription request
type UnsubscribeEvent struct {
	Type     string   `json:"type"`     // "UNSUBSCRIBE"
	Channels []string `json:"channels"` // List of channel names to unsubscribe from
	User     string   `json:"user"`     // Username (auto-filled)
}

// ListSubscriptionsEvent represents a request to list current subscriptions
type ListSubscriptionsEvent struct {
	Type string `json:"type"` // "LIST_SUBSCRIPTIONS"
	User string `json:"user"` // Username (auto-filled)
}

// SubscriptionResponseEvent represents the response to subscription requests
type SubscriptionResponseEvent struct {
	Type          string   `json:"type"`          // "SUBSCRIPTION_RESPONSE"
	Action        string   `json:"action"`        // "subscribed", "unsubscribed", "listed"
	Channels      []string `json:"channels"`      // Affected channels
	Subscriptions []string `json:"subscriptions"` // Current subscriptions (for list action)
	Success       bool     `json:"success"`
	Message       string   `json:"message,omitempty"`
}

// GetType returns the event type for SubscribeEvent
func (e *SubscribeEvent) GetType() string { return e.Type }

// GetUser returns the user for SubscribeEvent
func (e *SubscribeEvent) GetUser() string { return e.User }

// GetType returns the event type for UnsubscribeEvent
func (e *UnsubscribeEvent) GetType() string { return e.Type }

// GetUser returns the user for UnsubscribeEvent
func (e *UnsubscribeEvent) GetUser() string { return e.User }

// GetType returns the event type for ListSubscriptionsEvent
func (e *ListSubscriptionsEvent) GetType() string { return e.Type }

// GetUser returns the user for ListSubscriptionsEvent
func (e *ListSubscriptionsEvent) GetUser() string { return e.User }

// GetType returns the event type for SubscriptionResponseEvent
func (e *SubscriptionResponseEvent) GetType() string { return e.Type }

// GetUser returns empty string for response events
func (e *SubscriptionResponseEvent) GetUser() string { return "" }

// HandleSubscribe processes subscription requests
func (h *Handler) HandleSubscribe(client shared.ClientInterface, messageBytes []byte) error {
	var event SubscribeEvent
	if err := json.Unmarshal(messageBytes, &event); err != nil {
		return fmt.Errorf("invalid SUBSCRIBE event: %v", err)
	}

	if err := h.validator.ValidateSubscribe(&event); err != nil {
		return err
	}

	event.User = client.GetUsername()

	successChannels := []string{}
	failedChannels := []string{}

	// Subscribe to each requested channel
	for _, channelName := range event.Channels {
		if err := client.GetHub().SubscribeToChannel(client, channelName); err != nil {
			log.Printf("❌ Failed to subscribe %s to %s: %v", event.User, channelName, err)
			failedChannels = append(failedChannels, channelName)
		} else {
			successChannels = append(successChannels, channelName)
		}
	}

	// Send response
	response := &SubscriptionResponseEvent{
		Type:     "SUBSCRIPTION_RESPONSE",
		Action:   "subscribed",
		Channels: successChannels,
		Success:  len(successChannels) > 0,
	}

	if len(failedChannels) > 0 {
		response.Message = fmt.Sprintf("Failed to subscribe to: %v", failedChannels)
	}

	return client.GetHub().SendToClient(client, response)
}

// HandleUnsubscribe processes unsubscription requests
func (h *Handler) HandleUnsubscribe(client shared.ClientInterface, messageBytes []byte) error {
	var event UnsubscribeEvent
	if err := json.Unmarshal(messageBytes, &event); err != nil {
		return fmt.Errorf("invalid UNSUBSCRIBE event: %v", err)
	}

	if err := h.validator.ValidateUnsubscribe(&event); err != nil {
		return err
	}

	event.User = client.GetUsername()

	successChannels := []string{}

	// Unsubscribe from each requested channel
	for _, channelName := range event.Channels {
		if err := client.GetHub().UnsubscribeFromChannel(client, channelName); err != nil {
			log.Printf("❌ Failed to unsubscribe %s from %s: %v", event.User, channelName, err)
		} else {
			successChannels = append(successChannels, channelName)
		}
	}

	// Send response
	response := &SubscriptionResponseEvent{
		Type:     "SUBSCRIPTION_RESPONSE",
		Action:   "unsubscribed",
		Channels: successChannels,
		Success:  len(successChannels) > 0,
	}

	return client.GetHub().SendToClient(client, response)
}

// HandleListSubscriptions processes requests to list current subscriptions
func (h *Handler) HandleListSubscriptions(client shared.ClientInterface, messageBytes []byte) error {
	var event ListSubscriptionsEvent
	if err := json.Unmarshal(messageBytes, &event); err != nil {
		return fmt.Errorf("invalid LIST_SUBSCRIPTIONS event: %v", err)
	}

	event.User = client.GetUsername()

	// Get current subscriptions
	subscriptions := client.GetHub().GetClientSubscriptions(client)

	// Send response
	response := &SubscriptionResponseEvent{
		Type:          "SUBSCRIPTION_RESPONSE",
		Action:        "listed",
		Subscriptions: subscriptions,
		Success:       true,
	}

	return client.GetHub().SendToClient(client, response)
}
