package socket

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/auth/authid"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/ws/storage"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
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

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func readPump(c storage.Conn) {
	defer func() {
		if c.GetHub() != nil {
			c.GetHub().Unregister(c)
		}
		c.GetConn().Close()
	}()
	c.GetConn().SetReadLimit(maxMessageSize)
	c.GetConn().SetReadDeadline(time.Now().Add(pongWait))
	c.GetConn().SetPongHandler(func(string) error { c.GetConn().SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.GetConn().ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Warnf("readPump: %v", err)
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
func writePump(c storage.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.GetConn().Close()
	}()
	for {
		select {
		case message, ok := <-c.GetSendChan():
			c.GetConn().SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.GetConn().WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.GetConn().NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.GetSendChan())
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.GetSendChan())
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.GetConn().SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.GetConn().WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeUserWs(c *gin.Context) {
	w := c.Writer
	r := c.Request

	token := c.Query("token")
	var iden *authid.Identity
	if token != "" {
		idenFull, err := auth.ParseJwtToken(c.Request.Context(), token)
		if err == nil {
			iden = &idenFull.Identity
		}
	}

	if iden == nil {
		c.AbortWithStatus(401)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warnf("User ws upgrade error: %s", err)
		return
	}

	hub := storage.GetOrCreateUserHub(iden.UserId)
	uconn := storage.NewUserConn(iden, conn, hub)
	hub.Register(uconn)

	go readPump(uconn)
	go writePump(uconn)
}
