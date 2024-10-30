package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/table/services/model"
	"github.com/glossd/pokergloss/table/services/player"
	"net/http"
)

type SetIntentInput struct {
	Intent model.Intent `json:"intent"`
}

// @ID set intent
// @Param id path string true "Table ID"
// @Param position path int true "Seat position"
// @Param input body SetIntentInput true "input body"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /tables/{id}/seats/{position}/intent [put]
func PutIntent(c *gin.Context) {
	var input SetIntentInput
	err := c.BindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	posParams, err := PositionParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	params, err := player.NewIntentParams(posParams, input.Intent)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = player.SetIntent(*params)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, Ok())
}

// @ID delete intent
// @Param id path string true "Table ID"
// @Param position path int true "Seat position"
// @Success 200 {array} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /tables/{id}/seats/{position}/intent [delete]
func DeleteIntent(c *gin.Context) {
	params, err := PositionParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = player.RemoveIntent(params)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, Ok())
}
