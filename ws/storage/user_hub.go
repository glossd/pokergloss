package storage

import (
	"sync/atomic"
)

type UserHub struct {
	Broadcast chan []byte
	hub *UserHub
	register   chan *UserConn
	unregister chan *UserConn
	userConns  map[*UserConn]bool
}

func NewEmptyUserHub() *UserHub {
	return &UserHub{}
}

func (h *UserHub) InitEmptyHub() {
	h.Broadcast = make(chan []byte)
	h.register = make(chan *UserConn)
	h.unregister = make(chan *UserConn)
	h.userConns = make(map[*UserConn]bool)
}

func (h *UserHub) IsZeroUserConn() bool {
	if h == nil {
		return true
	}
	return len(h.userConns) == 0
}

func (h *UserHub) Register(conn Conn) {
	tconn, ok := conn.(*UserConn)
	if ok {
		h.register <- tconn
		atomic.AddInt64(&userConnsCount, 1)
	}
}

func (h *UserHub) Unregister(conn Conn) {
	tconn, ok := conn.(*UserConn)
	if ok {
		h.unregister <- tconn
		atomic.AddInt64(&userConnsCount, -1)
	}
}

func (h *UserHub) Run() {
	for {
		select {
		case tconn := <-h.register:
			h.userConns[tconn] = true
		case tconn := <-h.unregister:
			if _, ok := h.userConns[tconn]; ok {
				delete(h.userConns, tconn)
			}
		case message := <-h.Broadcast:
			for userConn := range h.userConns {
				select {
				case userConn.Send <- message:
				default:
					close(userConn.Send)
					delete(h.userConns, userConn)
				}
			}
		}
	}
}
