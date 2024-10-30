package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/market/service"
	"github.com/glossd/pokergloss/market/web/rest/model"
	"net/http"
)

// @ID buy item
// @Param input body model.BuyItemParams true "input body"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /items [post]
func BuyItem(c *gin.Context) {
	var input model.BuyItemParams
	err := c.BindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = service.BuyItem(c.Request.Context(), auth.Id(c), input)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, Ok("bought"))
}
