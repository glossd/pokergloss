package router

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/messenger/web/rest"
	"github.com/glossd/pokergloss/messenger/web/ws"
	"github.com/glossd/pokergloss/messenger/web/ws/wsstore"
	"net/http"
)

const BasePath = "/api/messenger"

func New(r *gin.Engine) *gin.Engine {
	base := r.Group(BasePath)
	base.GET("/status", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"up": 1}) })

	wsRouter := base.Group("/ws")
	wsRouter.Use(auth.WebsocketMiddleware)
	wsRouter.GET("/events", ws.ServeWs)

	authenticated := base.Group("/")
	authenticated.Use(auth.EmailVerifiedMiddleware)

	authenticated.GET("/chats", rest.GetAllChats)
	authenticated.POST("/u2u/chats", rest.PostU2UChat)

	authenticated.GET("/unread-chats-count", rest.GetUnreadChatsCount)

	authenticated.GET("/chats/:chatId/messages", rest.GetChatMessages)
	authenticated.POST("/chats/:chatId/messages", rest.SendChatMessage)
	authenticated.PUT("/chats/:chatId/read/messages", rest.ReadMessages)
	authenticated.GET("/use-ws", wsstore.UseWS)
	return r
}
