package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/market/service"
	"github.com/glossd/pokergloss/market/web/rest/model"
	"net/http"
)

// @ID list my purchases
// @Success 200 {array} model.Purchase
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /me/purchases [get]
func ListPurchases(c *gin.Context) {
	commands, err := service.ListBuyCommands(c.Request.Context(), auth.Id(c).UserId)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.ToPurchases(commands))
}
