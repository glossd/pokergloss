package e2e

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/cleaning"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestCleanFinishedLobbies(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.Table.GameEndMinTimeout = 0

	domain.Algo = Algo_2P_MockZeroPositionLoses(t)

	lobby := insertFullLobbySitAndGo(t, func() *domain.LobbySitAndGo { return initLobbySitngo(t) })
	table := lobby.GetTable()
	restMakeAction(t, table.ID.Hex(), domain.AllIn, 0)
	restMakeAction(t, table.ID.Hex(), domain.AllIn, 1, secondPlayerToken)

	deleteCount, err := cleaning.CleanFinishedLobbiesOf(db.ColSitAndGoLobby())
	assert.Nil(t, err)
	assert.EqualValues(t, 1, deleteCount)
}

func TestCleanNotStartedSitngo(t *testing.T) {
	prevPropsSetup(t)

	domain.Algo = Algo_2P_MockZeroPositionLoses(t)

	l, err := domain.NewLobbySitAndGo(domain.NewLobbySitAndGoParams{
		NewTableParams: domain.NewTableParams{
			Name:            "Test",
			Size:            6,
			BigBlind:        10,
			DecisionTimeout: domain.DefaultDecisionTimeout,
			BettingLimit:    domain.NL,
			Identity:        domain.SystemIdentity,
		},
		PlacesPaid:        1,
		BuyIn:             1000,
		LevelIncreaseTime: domain.DefaultLevelIncreaseTime,
	})
	assert.Nil(t, err)

	l.CreatedAt = 10
	err = db.InsertSitAndGoLobbyNoCtx(l)
	assert.Nil(t, err)

	cleaning.CleanNotStartedLobbies()
	_, err = db.FindSitAndGoLobbyNoCtx(l.ID)
	assert.EqualValues(t, mongo.ErrNoDocuments, err)
}
