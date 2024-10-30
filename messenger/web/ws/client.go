package ws

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/auth/authid"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/messenger/service"
	"github.com/glossd/pokergloss/messenger/web/ws/wsstore"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
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

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func readPump(c *wsstore.Client) {
	defer func() {
		wsstore.RemoveUser(c)
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var e ReadEvent
		err = json.Unmarshal(message, &e)
		if err != nil {
			log.Errorf("Failed to parse read event: %s", err)
		}

		ctx := context.Background()
		_ = service.HandleTyping(ctx, c.Identity, e.Payload.ChatID, e.Payload.Text)

		log.Debugf("Got message from ws: %s", message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func writePump(c *wsstore.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.SendChan:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			writeMessage(w, message)

			// Add queued chat messages to the current websocket message.
			n := len(c.SendChan)
			for i := 0; i < n; i++ {
				w.Write(newline)
				writeMessage(w, <-c.SendChan)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func writeMessage(w io.WriteCloser, event *wsstore.MessengerEvent) {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Errorf("")
	}
	w.Write(payload)
}

// ServeWs handles websocket requests from the peer.
func ServeWs(c *gin.Context) {
	w := c.Writer
	r := c.Request

	token := c.Query("token")
	var iden *authid.Identity
	if token != "" {
		fullIden, err := auth.ParseJwtToken(c.Request.Context(), token)
		if err != nil {
			log.Warn("couldn't parse jwt token")
			c.AbortWithStatus(401)
			return
		}
		iden = &fullIden.Identity
	}

	if iden == nil {
		c.AbortWithStatus(401)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn(err)
		return
	}
	client := wsstore.NewClient(*iden, conn)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go writePump(client)
	go readPump(client)
}
