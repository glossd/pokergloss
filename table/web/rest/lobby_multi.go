package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/services/model"
	"github.com/glossd/pokergloss/table/services/multi"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// @ID listMultiLobbies
// @Param skip query int false "results after a certain number"
// @Param limit query int false "the maximum number of results to be returned, max 20"
// @Param skipEmpty query bool false "skip empty"
// @Param skipFull query bool false "skip full"
// @Success 200 {array} model.LobbyMulti
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /multi/lobbies [get]
func GetMultiLobbies(c *gin.Context) {
	params, err := PageParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	lobbies, err := multi.FindLobbies(c.Request.Context(), *params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	c.JSON(http.StatusOK, model.ToMultiLobbies(lobbies))
}

// @ID registerMulti
// @Param id path string true "Lobby ID"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /multi/lobbies/{id}/register [put]
func RegisterMulti(c *gin.Context) {
	params, err := IdenParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	err = multi.Register(*params)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, Ok())
}

// @ID unregisterMulti
// @Param id path string true "Lobby ID"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /multi/lobbies/{id}/unregister [put]
func UnregisterMulti(c *gin.Context) {
	params, err := IdenParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	err = multi.Unregister(*params)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, Ok())
}

// @ID get full multi lobby
// @Param id path string true "Lobby ID"
// @Success 200 {object} model.LobbyMulti
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /multi/lobbies/{id} [get]
func GetFullLobbyMulti(c *gin.Context) {
	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(400, E(err))
		return
	}

	lobby, err := db.FindLobbyMultiNoCtx(oid)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	lobbyMulti := model.ToLobbyMulti(lobby)
	lobbyMulti.FillTables()
	c.JSON(http.StatusOK, lobbyMulti)
}
