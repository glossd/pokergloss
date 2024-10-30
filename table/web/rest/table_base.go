package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/model"
	"github.com/glossd/pokergloss/table/services/tables"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// @ID get tables
// @Param skip query int false "ratings to skip"
// @Param limit query int false "number of ratings to return, max 20"
// @Param skipEmpty query bool false "skip empty tables"
// @Param skipFull query bool false "skip full tables"
// @Param type query string false "Variants cashGame,sitAndGo,multi. Defaults to cashGame"
// @Success 200 {array} model.Table
// @Failure 400 {object} ErrorRes
// @Router /tables [get]
func GetTables(c *gin.Context) {
	params, err := PageParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	tableType := toTableType(c.Request.URL.Query().Get("type"))

	foundTables, err := tables.FindAll(c.Request.Context(), tableType, *params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	result := model.ToModelTables(foundTables)
	c.JSON(http.StatusOK, &result)
}

func toTableType(t string) domain.TableType {
	switch t {
	case string(domain.CashType), string(domain.SitngoType), string(domain.MultiType):
		return domain.TableType(t)
	default:
		return domain.CashType
	}
}

type PostTableInput struct {
	Name     string `json:"name"`
	Size     int    `json:"size"`
	BigBlind int64  `json:"bigBlind"`
	// Optional, defaults to 10 seconds
	DecisionTimeoutSec int                 `json:"decisionTimeoutSec"`
	BettingLimit       domain.BettingLimit `json:"bettingLimit"`
	IsPrivate          bool                `json:"isPrivate"`
}

// @ID post table
// @Param input body PostTableInput true "New table"
// @Success 200 {object} model.Table
// @Failure 400 {object} ErrorRes
// @Router /tables [post]
func PostTable(c *gin.Context) {
	var input PostTableInput
	err := c.BindJSON(&input)
	if err != nil {
		log.Warn("Couldn't create table: parse input: %", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	params, err := ToNewTableParams(input, auth.Id(c))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	table, err := tables.Create(c.Request.Context(), params)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// would be nice to return 201
	c.JSON(http.StatusOK, model.ToModelTable(table, model.ToPlayerInMultiLobby))
}

func ToNewTableParams(input PostTableInput, identity authid.Identity) (domain.NewTableParams, error) {
	return domain.NewTableParams{
		Name:            input.Name,
		Size:            input.Size,
		BigBlind:        input.BigBlind,
		DecisionTimeout: time.Duration(input.DecisionTimeoutSec) * time.Second,
		BettingLimit:    input.BettingLimit,
		IsPrivate:       input.IsPrivate,
		Identity:        identity,
	}, nil
}

// @ID get full table
// @Param id path string true "Table ID"
// @Success 200 {object} model.Table
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /tables/{id} [get]
func GetFullTable(c *gin.Context) {
	id := c.Param("id")
	t, err := tables.Find(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	iden, isIden := auth.IdSafe(c)
	var userID string
	if isIden {
		userID = iden.UserId
	}
	table := model.ToModelTableSeats(t, model.AllSeatsIdentifiedCards(t, userID))

	c.JSON(http.StatusOK, table)
}
