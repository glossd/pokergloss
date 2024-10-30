package router

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/market/web/rest"
	"net/http"
)

const BasePath = "/api/market"

func New(r *gin.Engine) *gin.Engine {
	base := r.Group(BasePath)

	base.GET("/status", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"up": 1}) })

	base.GET("/items", rest.ListItems)
	base.GET("/products", rest.ListProducts)
	base.GET("/users/:userId/items", rest.ListUserItems)

	authenticated := base.Group("/")
	authenticated.Use(auth.EmailVerifiedMiddleware)
	authenticated.POST("/items", rest.BuyItem)
	authenticated.POST("/products", rest.BuyProduct)
	authenticated.PUT("/items/:itemId/select", rest.SelectItem)
	authenticated.GET("/me/items", rest.ListMyItems)
	authenticated.GET("/me/purchases", rest.ListPurchases)
	return r
}
