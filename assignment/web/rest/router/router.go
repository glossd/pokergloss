package router

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/assignment/web/mq/ws"
	"github.com/glossd/pokergloss/assignment/web/rest"
	"github.com/glossd/pokergloss/auth"
	"net/http"
)

const BasePath = "/api/assignment"

func New(r *gin.Engine) *gin.Engine {
	base := r.Group(BasePath)
	base.GET("/status", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"up": 1}) })

	authenticated := base.Group("/")
	authenticated.Use(auth.Middleware)
	authenticated.GET("/my/daily/assignments", rest.GetMyDailyAssignments)
	authenticated.GET("/ws", ws.UseWS)

	return r
}
