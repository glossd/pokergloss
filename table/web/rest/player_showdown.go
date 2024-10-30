package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/player"
	"net/http"
)

// @ID make show down action
// @Param id path string true "Table ID"
// @Param position path int true "Seat position"
// @Param action path string true "Show Down Action type"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /tables/{id}/seats/{position}/showdown-actions/{action} [put]
func MakeShowDownAction(c *gin.Context) {
	params, err := PositionParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	actionType, err := toShowDownActionType(c.Param("action"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = player.MakeShowDownAction(params, actionType)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, Ok())
}

func toShowDownActionType(a string) (domain.ShowDownActionType, error) {
	switch a {
	case string(domain.Muck), string(domain.Show), string(domain.ShowLeft), string(domain.ShowRight):
		return domain.ShowDownActionType(a), nil
	}
	return "", fmt.Errorf("no such show down action type: %s", a)
}
