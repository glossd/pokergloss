package router

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/bonus/web/rest"
	"net/http"
)

const BasePath = "/api/bonus"

func New(r *gin.Engine) *gin.Engine {
	base := r.Group(BasePath)

	base.GET("/status", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"up": 1}) })

	authenticated := base.Group("/")
	authenticated.Use(auth.Middleware)
	authenticated.PUT("/daily-bonus", rest.TakeDailyBonus)
	return r
}
