package e2e

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/ws/storage"
	"github.com/glossd/pokergloss/ws/web/router"
	"github.com/gorilla/websocket"
	"github.com/pokerblow/go-httptestutil"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const defaultToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjUxMDM2YWYyZDgzOWE4NDJhZjQzY2VjZmJiZDU4YWYxYTc1OGVlYTIiLCJ0eXAiOiJKV1QifQ.eyJ1c2VybmFtZSI6ImRlbmlzIiwiaXNzIjoiaHR0cHM6Ly9zZWN1cmV0b2tlbi5nb29nbGUuY29tL3Bva2VyYmxvdyIsImF1ZCI6InBva2VyYmxvdyIsImF1dGhfdGltZSI6MTU5ODY0Njg3OCwidXNlcl9pZCI6IjJmUjVNWXlqcU1TTWd6THRka2RyS1prYzE4dDEiLCJzdWIiOiIyZlI1TVl5anFNU01nekx0ZGtkcktaa2MxOHQxIiwiaWF0IjoxNTk4NjQ2ODc4LCJleHAiOjE1OTg2NTA0NzgsImVtYWlsIjoiZGVuaXNnbG90b3Y5OEBtYWlsLnJ1IiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7ImVtYWlsIjpbImRlbmlzZ2xvdG92OThAbWFpbC5ydSJdfSwic2lnbl9pbl9wcm92aWRlciI6InBhc3N3b3JkIn19.Kcn9QroIR-62xGlzGTHdvx2uNRfiqAtUJyYBfzg74Mt_v4XozZTW-6O_teFmFoRasJsOr49uW4i9ntkkgoc6FgDoo1jTi_1yMMx3_gNS9qSAMMcmMscqelOHQdgxsi9mJMwltHqHOf-AsoYl7qbc_HCf5ShYBtljZlkUXY_pMGvy0ePupNMFiWxoTYmNiIelaz0d-O9oVzns8XOm6O6A5qsFLx6hnNNsS7cBMMbc9zqhsySAZhMzYHdkd-LvL8QtUUjAQQfsbit9hPFa4irFEf7gfYOX61kUeKlMBP6-f6O8Q8GmgieBJ30Ly8wQwJcgHcYkTAFlZhD8FRO08QesgQ"

var defaultIdentity = authid.Identity{UserId: "2fR5MYyjqMSMgzLtdkdrKZkc18t1", Username: "denis"}

var testRouter = httptestutil.NewRouter(router.New(gin.New())).BasePath(router.BasePath).Headers(authHeaders(defaultToken))

func authHeaders(token string) map[string]string {
	if token == "" {
		return map[string]string{"Authorization": ""}
	}
	return map[string]string{"Authorization": "Bearer " + token}
}

func evt(etype string) *mqws.Event {
	return &mqws.Event{
		Type:    etype,
		Payload: "{}",
	}
}

func cleanUp(t *testing.T) {
	t.Cleanup(func() {
		storage.ReInitStorage()
	})
}

func wsDial(t *testing.T, entityType mqws.Message_EntityType, token ...string) (ws *websocket.Conn, closeWS func()) {
	s := httptest.NewServer(router.New(gin.New()))
	ws, _, err := websocket.DefaultDialer.Dial(wsURL(s, entityType, tokenOrDefault(token)), nil)
	if err != nil {
		t.Fatal("WS Dial failed: ", err)
	}
	ws.SetReadDeadline(time.Now().Add(time.Second * 5))
	return ws, func() {
		s.Close()
		ws.Close()
	}
}

func tokenOrDefault(token []string) string {
	if len(token) > 0 {
		return token[0]
	}
	return defaultToken
}

func wsURL(s *httptest.Server, entityType mqws.Message_EntityType, token string) string {
	var restOfPath string
	switch entityType {
	case mqws.Message_USER:
		restOfPath = "/news"
	}
	url := fmt.Sprintf("ws://%s%s%s?&token=%s", strings.TrimPrefix(s.URL, "http://"), router.BasePath, restOfPath, token)
	return url
}
