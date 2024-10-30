package router

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/table-stream/web/rest"
	"github.com/glossd/pokergloss/table-stream/web/ws"
	"net/http"
)

const BasePath = "/api/table-stream"

func New(r *gin.Engine) *gin.Engine {
	base := r.Group(BasePath)

	// automatically add routers for net/http/pprof
	// e.g. /debug/pprof, /debug/pprof/heap, etc.
	ginpprof.WrapGroup(base)

	base.GET("/status", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"up": 1}) })

	base.GET("/tables/:id", ws.ServeWs)
	base.GET("/tables/:id/users", rest.GetTableUsers)

	return r
}
