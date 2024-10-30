package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/market/service"
	"net/http"
)

// @ID list products
// @Success 200 {array} model.Product
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /products [get]
func ListProducts(c *gin.Context) {
	c.JSON(http.StatusOK, service.ProductList)
}

type BuyProductParams struct {
	ID string `json:"id"`
}

// @ID buy product
// @Param input body BuyProductParams true "input body"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /products [post]
func BuyProduct(c *gin.Context) {
	var input BuyProductParams
	err := c.BindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = service.BuyProduct(c.Request.Context(), auth.Id(c), input.ID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, Ok("bought"))
}
