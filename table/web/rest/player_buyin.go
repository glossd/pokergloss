package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/table/services/player"
	"net/http"
)

type BuyInInputV2 struct {
	// amount of chips
	Stack int64 `json:"stack"`
}

// @ID buy in v2
// @Param id path string true "Table ID"
// @Param position path int true "Seat position"
// @Param input body BuyInInputV2 true "input body"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /tables/{id}/seats/{position}/buy-in [put]
func BuyInV2(c *gin.Context) {
	var input BuyInInputV2
	err := c.BindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	params, err := PositionChipsParams(c, input.Stack)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	_, err = player.BuyIn(params)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, Ok())
}
