package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/player"
	"net/http"
)

type ActionInputV2 struct {
	// Required only for bet and raise
	Chips int64 `json:"chips"`
}

// @ID make action v2
// @Param id path string true "Table ID"
// @Param position path int true "Seat position"
// @Param action path string true "Action type"
// @Param input body ActionInputV2 true "input body"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /tables/{id}/seats/{position}/actions/{action} [put]
func MakeActionV2(c *gin.Context) {
	var input ActionInputV2
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

	actionType, err := toActionType(c.Param("action"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = player.MakeBettingAction(params, actionType)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, Ok())
}

func toActionType(a string) (domain.ActionType, error) {
	switch a {
	case string(domain.Check), string(domain.Bet), string(domain.Fold),
		string(domain.Call), string(domain.Raise), string(domain.AllIn):
		return domain.ActionType(a), nil
	}
	return "", fmt.Errorf("no such action type")
}
