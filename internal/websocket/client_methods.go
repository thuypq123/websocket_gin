package websocket

import (
	"websocket/internal/websocket/handlers/shared"
)

// Make Client implement ClientInterface
func (c *Client) GetUsername() string {
	return c.username
}

func (c *Client) GetID() string {
	return c.id
}

func (c *Client) GetHub() shared.HubInterface {
	return c.hub
}

func (c *Client) SendError(message string) {
	errorEvent := shared.NewErrorEvent(message)
	c.hub.SendToClient(c, errorEvent)
}

// handleEvent routes events using the event router
func (c *Client) handleEvent(messageBytes []byte) error {
	return eventRouter.routeEvent(c, messageBytes)
}
