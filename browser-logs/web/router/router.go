package router

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/browser-logs/web/rest"
	"net/http"
)

const BasePath = "/api/browser-logs"

func New(r *gin.Engine) *gin.Engine {
	base := r.Group(BasePath)
	base.GET("/status", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"up": 1}) })
	base.POST("/errors", rest.PostError)
	return r
}
