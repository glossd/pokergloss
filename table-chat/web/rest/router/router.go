package router

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/table-chat/web/rest"
	"net/http"
)

const BasePath = "/api/table-chat"

func New(r *gin.Engine) *gin.Engine {
	base := r.Group(BasePath)
	base.GET("/status", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"up": 1}) })

	authenticated := base.Group("/")
	authenticated.Use(auth.EmailVerifiedMiddleware)
	authenticated.POST("/tables/:tableId/messages", rest.PostMessage)
	authenticated.POST("/tables/:tableId/emojis", rest.PostEmoji)
	return r
}
