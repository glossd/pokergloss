package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/table/services/player"
	"net/http"
)

type AddChipsInput struct {
	Chips int64 `json:"chips"`
}

// @ID add chips
// @Param id path string true "Table ID"
// @Param position path int true "Seat position"
// @Param input body AddChipsInput true "input body"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /tables/{id}/seats/{position}/add-chips [put]
func AddChips(c *gin.Context) {
	var input AddChipsInput
	err := c.BindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	params, err := PositionChipsParams(c, input.Chips)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = player.AddChips(params)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, Ok())
}
