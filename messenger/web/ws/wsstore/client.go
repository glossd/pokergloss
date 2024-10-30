package wsstore

import (
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Identity authid.Identity
	// The websocket connection.
	Conn *websocket.Conn

	SendChan chan *MessengerEvent
}

func NewClient(iden authid.Identity, conn *websocket.Conn) *Client {
	client := &Client{
		Identity: iden,
		Conn:     conn,
		SendChan: make(chan *MessengerEvent, 256),
	}
	addUser(client)
	return client
}

func (c *Client) Close() {
	close(c.SendChan)
	c.Conn.Close()
}
