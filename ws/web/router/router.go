package router

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/ws/web/rest"
	"github.com/glossd/pokergloss/ws/web/socket"
	"net/http"
)

const BasePath = "/api/ws"

func New(r *gin.Engine) *gin.Engine {
	base := r.Group(BasePath)

	// automatically add routers for net/http/pprof
	// e.g. /debug/pprof, /debug/pprof/heap, etc.
	ginpprof.WrapGroup(base)

	base.GET("/status", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"up": 1}) })

	base.GET("/news", socket.ServeUserWs)

	notification := base.Group("/notification")
	notification.Use(auth.Middleware)
	notification.PUT("/tokens", rest.CheckOrUpdateNotificationToken)
	return r
}
