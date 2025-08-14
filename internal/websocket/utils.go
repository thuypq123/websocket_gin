package websocket

import (
	"fmt"
	"time"
)

// generateClientID creates a unique client ID
func generateClientID() string {
	return fmt.Sprintf("client_%d", time.Now().UnixNano())
}

// GetStats returns simple hub statistics
func (h *Hub) GetStats() map[string]interface{} {
	h.roomsMutex.RLock()
	h.postMutex.RLock()
	defer h.roomsMutex.RUnlock()
	defer h.postMutex.RUnlock()

	return map[string]interface{}{
		"total_clients":    len(h.clients),
		"chat_rooms":       len(h.chatRooms),
		"post_subscribers": len(h.postSubscribers),
	}
}
