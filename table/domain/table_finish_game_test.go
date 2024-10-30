package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTable_Reset(t *testing.T) {
	table := table2PlayerPositions_gameRiver(t, 0, 1)

	table.reset()

	assert.EqualValues(t, 1, len(table.Pots))
	assert.EqualValues(t, 0, table.Pots[0].Chips)
	assert.EqualValues(t, 0, table.TotalPot)
	assert.EqualValues(t, CommunityCards{}, *table.CommunityCards)
	assert.EqualValues(t, -1, table.DecidingPosition)
	assert.Nil(t, table.Winners)
}

func Test_SbFold_CommunityCardsNotFull(t *testing.T) {
	table := table2PlayersPositions_startedGame(t, 0, 1)

	err := table.MakeActionDeprecated(table.DecidingPosition, Fold,0)
	assert.Nil(t, err)

	assert.Len(t, table.CommunityCards.Available(), 0)
}

func Test_Blind_WithNotEnoughChips(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table2PlayersPositions_startedGame(t, 0, 1)
	sb := table.DecidingPlayerUnsafe()
	err := table.MakeActionDeprecated(sb.Position, Raise, 248)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, sb.Stack)

	bb := table.DecidingPlayerUnsafe()
	err = table.MakeActionDeprecated(bb.Position, Call, 0)
	assert.Nil(t, err)

	err = table.MakeActionDeprecated(bb.Position, Fold, 0)
	assert.Nil(t, err)


	brokeP := bb
	// broke player with one chip should lose
	mock, err := NewMockAlgo(CardsStr("As", "Ks", "2d", "7h", "Qs", "Js", "Ts", "8h", "4c"))
	assert.Nil(t, err)
	Algo = mock

	assert.Nil(t, table.StartNextGame())

	assert.Zero(t, brokeP.Stack)
	assert.EqualValues(t, DealerSmallBlind, brokeP.Blind)
	assert.EqualValues(t, 1, brokeP.TotalRoundBet)
	assert.EqualValues(t, 1, brokeP.LastGameBet)
	assert.EqualValues(t, AllIn, brokeP.LastRoundAction)
	assert.EqualValues(t, AllIn, brokeP.LastGameAction)
	// assert.EqualValues(t, PlayerSittingOut, brokeP.Status) should sit out right away?

	assert.EqualValues(t, -1, table.DecidingPosition)
	assert.True(t, table.IsGameEnd())

	assert.Nil(t, table.StartNextGame())
	assert.True(t, table.IsWaiting())
}