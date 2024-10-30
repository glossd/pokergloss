package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTable_Stand_AfterReserve(t *testing.T) {
	table := defaultTable(t)
	assert.Nil(t, table.ReserveSeat(0, firstIdentity))

	_, err := table.Stand(0, firstIdentity)
	assert.Nil(t, err)
}

func TestTable_Stand_AfterBuyIn(t *testing.T) {
	table := defaultTable(t)
	assert.Nil(t, table.ReserveSeat(0, firstIdentity))
	_, err := table.BuyIn(defaultBuyIn, 0, firstIdentity)
	assert.Nil(t, err)

	_, err = table.Stand(0, firstIdentity)
	assert.Nil(t, err)
}

func TestTable_Stand_WhilePlaying(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table2PlayersPositions_startedGame(t, 0, 1)

	_, err := table.Stand(1, secondIdentity)
	assert.Nil(t, err)

	assert.True(t, table.IsInPlay())

	table.tcall(t)

	assert.True(t, table.IsGameEnd())
}

func TestTable_Stand_WhileDeciding_ShouldFold(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table2PlayersPositions_startedGame(t, 0, 1)

	table.tcall(t)
	_, err := table.Stand(1, secondIdentity)
	assert.Nil(t, err)

	assert.True(t, table.IsGameEnd())
}

func TestTable_Stand_ShouldNotRemoveSeatBlind(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table3Players_startedGame(t)

	table.tcall(t)
	sbPos := table.DecidingPosition
	_, err := table.Stand(sbPos, thirdIdentity)
	assert.Nil(t, err)

	assert.True(t, table.IsSeatFree(sbPos))
	assert.EqualValues(t, sbPos, table.SmallBlindPosition())

	table.tfold(t)
	assert.True(t, table.IsGameEnd())

	assert.EqualValues(t, SmallBlind, table.GetSeatUnsafe(sbPos).Blind)
	assert.Nil(t, table.StartNextGame())
	assert.EqualValues(t, "", table.GetSeatUnsafe(sbPos).Blind)
}

func TestTable_Stand_ShouldStandAfterFold(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table3Players_startedGame(t)

	table.tcall(t)
	table.tfold(t)
	table.tcheck(t)

	_, err := table.Stand(table.SmallBlindPosition(), thirdIdentity)
	assert.Nil(t, err)

	assert.True(t, table.IsSeatFree(table.SmallBlindPosition()))
}
