package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/model"
	"github.com/glossd/pokergloss/table/services/sitngo"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type PostSitAndGoLobbyInput struct {
	TableParams              PostTableInput      `json:"tableParams"`
	PlacesPaid               int                 `json:"placesPaid"`
	BuyIn                    int64               `json:"buyIn"`
	LevelIncreaseTimeMinutes int                 `json:"levelIncreaseTimeMinutes"`
	BettingLimit             domain.BettingLimit `json:"bettingLimit"`
	StartAt                  int64               `json:"startAt"`
}

// @ID postSitAndGoLobby
// @Param input body PostSitAndGoLobbyInput true "Lobby params"
// @Success 201 {object} model.LobbySitAndGo
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /sit-n-go/lobbies [post]
func PostSitAndGoLobby(c *gin.Context) {
	var input PostSitAndGoLobbyInput
	err := c.BindJSON(&input)
	if err != nil {
		log.Warn("Couldn't create sig&go lobby: parse input: %", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	params, err := toNewSitAndGoLobbyParams(input, auth.Id(c))
	if err != nil {
		log.Warn("Couldn't create sig&go lobby: creating params: %", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	lobby, err := sitngo.Create(c.Request.Context(), params)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, lobby)
}

func toNewSitAndGoLobbyParams(input PostSitAndGoLobbyInput, iden authid.Identity) (domain.NewLobbySitAndGoParams, error) {
	tableParams, err := ToNewTableParams(input.TableParams, iden)
	if err != nil {
		return domain.NewLobbySitAndGoParams{}, err
	}

	return domain.NewLobbySitAndGoParams{
		NewTableParams:    tableParams,
		PlacesPaid:        input.PlacesPaid,
		BuyIn:             input.BuyIn,
		LevelIncreaseTime: time.Duration(input.LevelIncreaseTimeMinutes) * time.Minute,
		StartAt:           input.StartAt,
	}, nil
}

// @ID listSitAndGoLobbies
// @Param skip query int false "results after a certain number"
// @Param limit query int false "the maximum number of results to be returned, max 20"
// @Param skipEmpty query bool false "skip empty"
// @Param skipFull query bool false "skip full"
// @Success 200 {array} model.LobbySitAndGo
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /sit-n-go/lobbies [get]
func GetSitAndGoLobbies(c *gin.Context) {
	params, err := PageParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	lobbies, err := sitngo.FindAll(c.Request.Context(), *params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	c.JSON(http.StatusOK, model.ToSNGLobbies(lobbies))
}

type RegisterForSitAndGoInput struct {
	Position int `json:"position"`
}

// @ID registerForSitAndGo
// @Param id path string true "Lobby ID"
// @Param input body RegisterForSitAndGoInput true "Register params"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /sit-n-go/lobbies/{id}/register [put]
func RegisterForSitAndGo(c *gin.Context) {
	var input RegisterForSitAndGoInput
	err := c.BindJSON(&input)
	if err != nil {
		log.Warn("Couldn't register sig&go lobby: parse input: %", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = sitngo.Register(c.Request.Context(), c.Param("id"), input.Position, auth.Id(c))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, Ok())
}

type UnregisterForSitAndGoInput struct {
	Position int `json:"position"`
}

// @ID unregisterFromSitAndGo
// @Param id path string true "Lobby ID"
// @Param input body UnregisterForSitAndGoInput true "Unregister params"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /sit-n-go/lobbies/{id}/unregister [put]
func UnregisterFromSitAndGo(c *gin.Context) {
	var input UnregisterForSitAndGoInput
	err := c.BindJSON(&input)
	if err != nil {
		log.Warn("Couldn't register sig&go lobby: parse input: %", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = sitngo.Unregister(c.Request.Context(), c.Param("id"), input.Position, auth.Id(c))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, Ok())
}

// @ID get full sitngo lobby
// @Param id path string true "Lobby ID"
// @Success 200 {object} model.LobbySitAndGo
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /sit-n-go/lobbies/{id} [get]
func GetFullSitngo(c *gin.Context) {
	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(400, E(err))
		return
	}

	lobby, err := db.FindSitAndGoLobbyNoCtx(oid)
	if err != nil {
		c.AbortWithStatusJSON(500, E(err))
		return
	}
	c.JSON(http.StatusOK, model.ToSitAndGoLobby(lobby))
}
