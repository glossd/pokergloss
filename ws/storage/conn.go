package storage

import (
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/gorilla/websocket"
)

type Conn interface {
	GetConn() *websocket.Conn
	GetSendChan() chan []byte
	GetHub() Hub
	GetIden() *authid.Identity
}

type UserConn struct {
	hub  *UserHub
	iden *authid.Identity
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	Send chan []byte
}

func NewUserConn(iden *authid.Identity, conn *websocket.Conn, hub *UserHub) *UserConn {
	return &UserConn{
		hub:  hub,
		iden: iden,
		conn: conn,
		Send: make(chan []byte, 256),
	}
}

func (u *UserConn) GetConn() *websocket.Conn {
	return u.conn
}

func (u *UserConn) GetSendChan() chan []byte {
	return u.Send
}

func (u *UserConn) GetHub() Hub {
	return u.hub
}

func (u *UserConn) GetIden() *authid.Identity {
	return u.iden
}
