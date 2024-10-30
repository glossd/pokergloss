package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/auth/authid"
	conf "github.com/glossd/pokergloss/goconf"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header["Origin"]
		if len(origin) == 0 {
			return true
		}
		u, err := url.Parse(origin[0])
		if err != nil {
			return false
		}

		host := u.Hostname()
		if strings.HasSuffix(host, conf.Props.Domain) {
			return true
		}
		if host == "localhost" || host == "127.0.0.1" {
			return true
		}

		if conf.IsDev() {
			return true
		}

		return false
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	iden *authid.Identity

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		atomic.AddInt64(&tableConnsCount, -1)
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Warnf("error: %v", err)
			}
			break
		}
		log.Debugf("Got message from ws: %s", message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(c *gin.Context) {
	w := c.Writer
	r := c.Request
	tableID := c.Param("id")
	token := c.Query("token")
	var iden *authid.Identity
	if token != "" {
		idenFull, err := auth.ParseJwtToken(c.Request.Context(), token)
		if err == nil {
			iden = &idenFull.Identity
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warnf("Table ws upgrade error: %s", err)
		return
	}

	var hub = getOrCreateTableHub(tableID)
	client := &Client{hub: hub, iden: iden, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	atomic.AddInt64(&tableConnsCount, 1)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
