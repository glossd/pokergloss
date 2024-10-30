package router

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/achievement/web/rest"
	"github.com/glossd/pokergloss/auth"
	"net/http"
)

const BasePath = "/api/achievement"

func New(r *gin.Engine) *gin.Engine {
	base := r.Group(BasePath)

	base.GET("/status", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"up": 1}) })

	base.GET("/users/:userId/points", rest.GetUserExp)
	base.GET("/users/:userId/achievements", rest.GetUserAchievements)
	base.GET("/users/:userId/statistics", rest.GetUserStatistics)

	authenticated := base.Group("/")
	authenticated.Use(auth.Middleware)
	authenticated.GET("/points/me", rest.GetMyExp)
	authenticated.GET("/achievements/me", rest.GetMyAchievements)
	authenticated.GET("/statistics/me", rest.GetMyStatistics)
	return r
}
