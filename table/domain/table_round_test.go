package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_2P_NewRounds(t *testing.T) {
	table := table2PlayersPositions_startedGame(t, 0, 1)
	assert.True(t, table.IsNewRound()) // preFlop is a new round, too

	err := table.MakeActionDeprecated(table.DecidingPosition, Call, 0)
	assert.Nil(t, err)
	assert.False(t, table.IsNewRound())
	err = table.MakeActionDeprecated(table.DecidingPosition, Check, 0)
	assert.Nil(t, err)
	assert.True(t, table.IsNewRound())
	assert.True(t, table.IsFlop())

	err = table.MakeActionDeprecated(table.DecidingPosition, Check, 0)
	assert.Nil(t, err)
	assert.False(t, table.IsNewRound())
	err = table.MakeActionDeprecated(table.DecidingPosition, Check, 0)
	assert.Nil(t, err)
	assert.True(t, table.IsNewRound())
	assert.True(t, table.IsTurn())

	err = table.MakeActionDeprecated(table.DecidingPosition, Check, 0)
	assert.Nil(t, err)
	assert.False(t, table.IsNewRound())
	err = table.MakeActionDeprecated(table.DecidingPosition, Check, 0)
	assert.Nil(t, err)
	assert.True(t, table.IsNewRound())
	assert.True(t, table.IsRiver())
}

func Test_2P_AfterPreFlop_LastShouldTakeActionDealer(t *testing.T) {
	table := table2PlayersPositions_startedGame(t, 0, 1)

	dealerSb := table.DecidingPlayerUnsafe()
	assert.EqualValues(t, DealerSmallBlind, dealerSb.Blind)

	err := table.MakeActionDeprecated(dealerSb.Position, Call, 0)
	assert.Nil(t, err)

	bb := table.DecidingPlayerUnsafe()
	assert.EqualValues(t, BigBlind, bb.Blind)

	err = table.MakeActionDeprecated(bb.Position, Check, 0)
	assert.Nil(t, err)


	assert.EqualValues(t, bb.Position, table.DecidingPosition)
}
