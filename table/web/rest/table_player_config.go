package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/table/services/model"
	"github.com/glossd/pokergloss/table/services/player"
	"net/http"
)

type AutoMuckConfig struct {
	AutoMuck bool `json:"autoMuck"`
}

// @ID set auto muck
// @Param id path string true "Table ID"
// @Param position path int true "Seat position"
// @Param autoMuckConfig body AutoMuckConfig true "Player's config"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /tables/{id}/seats/{position}/configs/auto-muck [put]
func SetAutoMuckPlayerConfig(c *gin.Context) {
	params, err := PositionParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	var input AutoMuckConfig
	err = c.BindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = player.SetAutoMuck(params, input.AutoMuck)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, Ok())
}

type AutoTopUpConfig struct {
	AutoTopUp bool `json:"autoTopUp"`
}

// @ID set auto top up
// @Param id path string true "Table ID"
// @Param position path int true "Seat position"
// @Param autoConfig body AutoTopUpConfig true "Player's config"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /tables/{id}/seats/{position}/configs/auto-top-up [put]
func SetAutoTopUpPlayerConfig(c *gin.Context) {
	params, err := PositionParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	var input AutoTopUpConfig
	err = c.BindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = player.SetAutoTopUp(params, input.AutoTopUp)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, Ok())
}

type AutoRebuyConfig struct {
	AutoReBuy bool `json:"autoRebuy"`
}

// @ID set auto rebuy
// @Param id path string true "Table ID"
// @Param position path int true "Seat position"
// @Param autoConfig body AutoRebuyConfig true "Player's config"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /tables/{id}/seats/{position}/configs/auto-rebuy [put]
func SetAutoRebuyPlayerConfig(c *gin.Context) {
	params, err := PositionParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	var input AutoRebuyConfig
	err = c.BindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = player.SetAutoReBuy(params, input.AutoReBuy)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, Ok())
}

// @ID get player config
// @Param id path string true "Table ID"
// @Param position path int true "Seat position"
// @Success 200 {object} model.PlayerAutoConfig
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /tables/{id}/seats/{position}/config [get]
func GetConfig(c *gin.Context) {
	params, err := PositionParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	config, err := player.GetConfig(params)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.ToPlayerAutoConfig(config))
}
