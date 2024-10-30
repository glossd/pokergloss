package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/market/domain"
	"github.com/glossd/pokergloss/market/web/rest/model"
	"net/http"
)

// @ID list items
// @Success 200 {array} model.Item
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /items [get]
func ListItems(c *gin.Context) {
	c.JSON(http.StatusOK, model.ToItems(domain.ItemsOnSale))
}
