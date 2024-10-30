package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/table/services/player"
	"net/http"
)

// @ID stand
// @Param id path string true "Table ID"
// @Param position path int true "Seat position"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /tables/{id}/seats/{position}/stand [delete]
func Stand(c *gin.Context) {
	params, err := PositionParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = player.Stand(params)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(200, Ok())
}
