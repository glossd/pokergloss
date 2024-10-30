package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/player"
	"net/http"
)

type PlayerStatusUpdate struct {
	Message PlayerStatus `json:"message"`
}

type PlayerStatus struct {
	Status domain.PlayerStatus `json:"status"`
}

// @ID back to game
// @Param id path string true "Table ID"
// @Param position path int true "Seat position"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /tables/{id}/seats/{position}/sit-back [put]
func SitBack(c *gin.Context) {
	params, err := PositionParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	p, err := player.SitBack(params)
	if err == domain.ErrNoSitBackInSittingOut && p != nil {
		res := PlayerStatusUpdate{Message: PlayerStatus{Status: p.Status}}
		c.AbortWithStatusJSON(http.StatusConflict, res)
		return
	}
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, Ok())
}
