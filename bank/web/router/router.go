package router

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/bank/web/rest"
	"net/http"
)

const BasePath = "/api/bank"

func New(r *gin.Engine) *gin.Engine {
	base := r.Group(BasePath)

	base.GET("/status", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"up": 1}) })

	base.GET("/balances/:userId", rest.GetUserBalanceWithRank)
	base.GET("/ratings/page", rest.GetRatings)

	authenticated := base.Group("/")
	authenticated.Use(auth.Middleware)
	authenticated.GET("/balance", rest.GetBalance)
	authenticated.GET("/ratings/me", rest.GetUserRating)

	authenticated.GET("/operations", rest.GetOperations)
	return r
}
