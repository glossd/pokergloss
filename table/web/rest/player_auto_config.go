package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/table/services/model"
	"github.com/glossd/pokergloss/table/services/playerautoconfig"
	"net/http"
)

// @ID set player auto config
// @Param config body model.PlayerAutoConfig true "Player's config"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /configs [put]
func SetPlayerAutoConfig(c *gin.Context) {
	var input model.PlayerAutoConfig
	err := c.BindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = playerautoconfig.SetAll(c.Request.Context(), auth.Id(c), &input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, Ok())
}

// @ID get player auto config
// @Success 200 {object} model.PlayerAutoConfig
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /configs [get]
func GetPlayerAutoConfig(c *gin.Context) {
	pac, err := playerautoconfig.FindAutoConfig(c.Request.Context(), auth.Id(c))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.ToPlayerAutoConfig(pac))
}
