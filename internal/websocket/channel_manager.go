package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"websocket/internal/websocket/handlers/shared"
)

// ChannelType represents different types of channels
type ChannelType string

const (
	ChannelTypeChat    ChannelType = "chat"
	ChannelTypePost    ChannelType = "post"
	ChannelTypeUser    ChannelType = "user"
	ChannelTypeSystem  ChannelType = "system"
	ChannelTypePrivate ChannelType = "private"
)

// Channel represents a communication channel
type Channel struct {
	Name        string      `json:"name"`
	Type        ChannelType `json:"type"`
	Identifier  string      `json:"identifier"`
	Description string      `json:"description,omitempty"`
	IsPrivate   bool        `json:"is_private"`
}

// ParseChannelName parses a channel name into type and identifier
func ParseChannelName(channelName string) (*Channel, error) {
	parts := strings.SplitN(channelName, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid channel name format: %s", channelName)
	}

	channelType := ChannelType(parts[0])
	identifier := parts[1]

	return &Channel{
		Name:       channelName,
		Type:       channelType,
		Identifier: identifier,
		IsPrivate:  channelType == ChannelTypeUser || channelType == ChannelTypePrivate,
	}, nil
}

// SubscribeToChannel subscribes a client to a channel
func (h *Hub) SubscribeToChannel(client shared.ClientInterface, channelName string) error {
	// Parse and validate channel
	channel, err := ParseChannelName(channelName)
	if err != nil {
		return err
	}

	// Validate permissions
	if err := h.validateChannelAccess(client, channel); err != nil {
		return err
	}

	concreteClient, ok := client.(*Client)
	if !ok {
		return fmt.Errorf("invalid client type")
	}

	// Add to hub's channel map
	h.channelsMutex.Lock()
	if h.channels[channelName] == nil {
		h.channels[channelName] = make(map[*Client]bool)
	}
	h.channels[channelName][concreteClient] = true
	h.channelsMutex.Unlock()

	// Add to client's subscription list
	concreteClient.subsMutex.Lock()
	if concreteClient.subscriptions == nil {
		concreteClient.subscriptions = make(map[string]bool)
	}
	concreteClient.subscriptions[channelName] = true
	concreteClient.subsMutex.Unlock()

	log.Printf("📡 Client %s subscribed to channel: %s", client.GetUsername(), channelName)
	return nil
}

// UnsubscribeFromChannel unsubscribes a client from a channel
func (h *Hub) UnsubscribeFromChannel(client shared.ClientInterface, channelName string) error {
	concreteClient, ok := client.(*Client)
	if !ok {
		return fmt.Errorf("invalid client type")
	}

	// Remove from hub's channel map
	h.channelsMutex.Lock()
	if channelClients, exists := h.channels[channelName]; exists {
		delete(channelClients, concreteClient)
		// Clean up empty channels
		if len(channelClients) == 0 {
			delete(h.channels, channelName)
		}
	}
	h.channelsMutex.Unlock()

	// Remove from client's subscription list
	concreteClient.subsMutex.Lock()
	delete(concreteClient.subscriptions, channelName)
	concreteClient.subsMutex.Unlock()

	log.Printf("📡 Client %s unsubscribed from channel: %s", client.GetUsername(), channelName)
	return nil
}

// BroadcastToChannel broadcasts a message to all subscribers of a channel
func (h *Hub) BroadcastToChannel(channelName string, event interface{}) error {
	h.channelsMutex.RLock()
	channelClients := h.channels[channelName]
	h.channelsMutex.RUnlock()

	if len(channelClients) == 0 {
		log.Printf("⚠️ No subscribers for channel: %s", channelName)
		return nil
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshaling event: %v", err)
	}

	successCount := 0
	for client := range channelClients {
		select {
		case client.send <- eventBytes:
			successCount++
		default:
			// Client's send buffer is full, remove it
			delete(channelClients, client)
			close(client.send)
			log.Printf("⚠️ Removed unresponsive client from channel: %s", channelName)
		}
	}

	log.Printf("📡 Broadcasted to channel %s (%d/%d clients)", channelName, successCount, len(channelClients))
	return nil
}

// GetClientSubscriptions returns all channels a client is subscribed to
func (h *Hub) GetClientSubscriptions(client shared.ClientInterface) []string {
	concreteClient, ok := client.(*Client)
	if !ok {
		return []string{}
	}

	concreteClient.subsMutex.RLock()
	defer concreteClient.subsMutex.RUnlock()

	subscriptions := make([]string, 0, len(concreteClient.subscriptions))
	for channel := range concreteClient.subscriptions {
		subscriptions = append(subscriptions, channel)
	}

	return subscriptions
}

// SendToUser sends a message to a specific user by username
func (h *Hub) SendToUser(username string, event interface{}) error {
	h.usersMutex.RLock()
	client := h.userConnections[username]
	h.usersMutex.RUnlock()

	if client == nil {
		return fmt.Errorf("user %s not connected", username)
	}

	return h.SendToClient(client, event)
}

// GetChannelSubscribers returns the number of subscribers for a channel
func (h *Hub) GetChannelSubscribers(channelName string) int {
	h.channelsMutex.RLock()
	defer h.channelsMutex.RUnlock()

	if channelClients, exists := h.channels[channelName]; exists {
		return len(channelClients)
	}
	return 0
}

// GetAllChannels returns all active channels
func (h *Hub) GetAllChannels() []string {
	h.channelsMutex.RLock()
	defer h.channelsMutex.RUnlock()

	channels := make([]string, 0, len(h.channels))
	for channelName := range h.channels {
		channels = append(channels, channelName)
	}

	return channels
}

// validateChannelAccess validates if a client can access a channel
func (h *Hub) validateChannelAccess(client shared.ClientInterface, channel *Channel) error {
	switch channel.Type {
	case ChannelTypeUser:
		// User channels are private to the user
		if channel.Identifier != client.GetUsername() {
			return fmt.Errorf("access denied to user channel: %s", channel.Name)
		}
	case ChannelTypePrivate:
		// Private channels need permission validation
		if !h.canAccessPrivateChannel(client, channel.Identifier) {
			return fmt.Errorf("access denied to private channel: %s", channel.Name)
		}
	case ChannelTypeChat, ChannelTypePost, ChannelTypeSystem:
		// Public channels - anyone can subscribe
		return nil
	default:
		return fmt.Errorf("unknown channel type: %s", channel.Type)
	}

	return nil
}

// canAccessPrivateChannel checks if user can access a private channel
func (h *Hub) canAccessPrivateChannel(client shared.ClientInterface, identifier string) bool {
	// For private conversations, check if user is participant
	// identifier format: "user1-user2" or conversation_id
	username := client.GetUsername()
	return strings.Contains(identifier, username)
}
