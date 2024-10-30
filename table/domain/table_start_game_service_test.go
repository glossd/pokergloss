package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_2P_StartNextGame(t *testing.T) {
	table := table2PlayerPositions_gameRiver(t, 0, 1)
	bbPosition := table.DecidingPosition
	err := table.MakeActionDeprecated(bbPosition, Check, 0)
	assert.Nil(t, err)
	dealerSbPosition := table.DecidingPosition
	err = table.MakeActionDeprecated(dealerSbPosition, Check, 0)
	assert.Nil(t, err)

	assert.True(t, table.IsGameEnd())

	assert.Nil(t, table.StartNextGame())

	assert.True(t, table.IsNewRound())
	assert.True(t, table.IsPreFlop())
	assert.True(t, table.IsNewGameStarted())

	assert.EqualValues(t, 1, len(table.Pots))
	assert.Zero(t, table.Pots[0].Chips)

	assert.NotNil(t, table.CommunityCards)
	assert.Nil(t, table.CommunityCards.Flop)
	assert.Nil(t, table.CommunityCards.Turn)
	assert.Nil(t, table.CommunityCards.River)

	assert.EqualValues(t, bbPosition, table.DecidingPosition)
	assert.EqualValues(t, DealerSmallBlind, table.DecidingPlayerUnsafe().Blind)
}

func TestTable_StopTable_AfterActionTimeout(t *testing.T) {
	table := table2PlayersPositions_startedGame(t, 0, 1)
	sbPosition := table.DecidingPosition

	assert.Nil(t,  table.MakeActionDeprecated(sbPosition, Call, 0))
	assert.Nil(t, table.MakeActionOnTimeout(table.DecidingPosition))
	assert.Nil(t, table.StartNextGame())
	assert.True(t, table.IsWaiting())
}
