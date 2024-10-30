package router

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/survival/web/rest"
	"net/http"
)

const BasePath = "/api/survival"

func New(r *gin.Engine) *gin.Engine {
	base := r.Group(BasePath)
	base.GET("/status", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"up": 1}) })
	base.GET("/users/:userId/score", rest.GetUserScore)

	anonymous := base.Group("/anonymous")
	anonymous.Use(auth.MiddlewareAnonymous)
	anonymous.POST("/start", rest.StartSurvivalAnonymous)

	authenticated := base.Group("/")
	authenticated.Use(auth.EmailVerifiedMiddleware)
	authenticated.POST("/start", rest.StartSurvival)
	authenticated.POST("/start-idle", rest.StartSurvivalIdle)
	authenticated.DELETE("/release", rest.ReleaseSurvival)
	authenticated.GET("/my/tickets", rest.GetTickets)
	return r
}
