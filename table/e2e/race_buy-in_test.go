package e2e

import (
	"context"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/player"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestBuyInStartGameRace(t *testing.T) {
	t.Cleanup(cleanUp)

	table := InsertTable(t)

	restReserveSeat(t, table.ID.Hex(), 0, getToken(0))
	restReserveSeat(t, table.ID.Hex(), 1, getToken(1))

	table1 := findTable(t, table.ID)
	table2 := findTable(t, table.ID)

	isNewGame, err := player.BuyInOnTable(table1, chipsParams(t, table.ID, 0))
	assert.False(t, isNewGame)
	assert.Nil(t, err)
	isNewGame, err = player.BuyInOnTable(table2, chipsParams(t, table.ID, 1))
	assert.False(t, isNewGame)
	assert.Equal(t, player.ErrBuyInRaceCondition, err)
}

func TestBuyIn_GameEndRace(t *testing.T) {
	t.Cleanup(cleanUp)
	conf.Props.GameEndMinTimeout = -1

	tableID := RestCreatedTableWithStartedGame(t, 0, 1)
	restReserveSeat(t, tableID.Hex(), 4, getToken(4))
	restMakeAction(t, tableID.Hex(), domain.Call, 0, getToken(0))
	restPlayerStand(t, tableID, 1, getToken(1))
	table := findTable(t, tableID)
	assert.EqualValues(t, domain.GameEndTable, table.Status)

	restPlayerStand(t, tableID, 4, getToken(4))
	tableStand := findTable(t, tableID)
	assert.EqualValues(t, table.GameFlowVersion, tableStand.GameFlowVersion)

	restReserveSeat(t, tableID.Hex(), 1, getToken(1))
	restBuyIn(t, tableID.Hex(), 1, getToken(1))
	table2 := findTable(t, tableID)
	assert.EqualValues(t, table.GameFlowVersion, table2.GameFlowVersion)

	restReserveSeat(t, tableID.Hex(), 2, getToken(2))
	restBuyIn(t, tableID.Hex(), 2, getToken(2))
	table3 := findTable(t, tableID)
	assert.EqualValues(t, table.GameFlowVersion, table3.GameFlowVersion)
}

func chipsParams(t *testing.T, tableID primitive.ObjectID, pos int) *player.ChipsParams {
	posParams1, err := player.NewPositionParams(context.Background(), tableID.Hex(), pos, getIden(pos))
	assert.Nil(t, err)
	params1, err := player.ToChipsParams(posParams1, 150)
	assert.Nil(t, err)
	return params1
}
