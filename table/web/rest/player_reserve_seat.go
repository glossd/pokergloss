package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/table/services/player"
	"net/http"
)

// @ID reserve seat v2
// @Param id path string true "Table ID"
// @Param position path int true "Seat position"
// @Success 200 {array} OkRes
// @Failure 400 {object} ErrorRes
// @Router /tables/{id}/seats/{position}/reserve [post]
func ReserveSeat(c *gin.Context) {
	params, err := PositionParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = player.ReserveTableSeat(params)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, Ok())
}
