package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReject_RaiseWithNotEnoughChips(t *testing.T) {
	table := table2PlayersPositions_startedGame(t, 0, 1)
	err := table.MakeActionDeprecated(table.DecidingPosition, Raise, 1)
	assert.Equal(t, ErrRaiseMoreChips, err)
}

// Game Cases

func TestReject_SbCheckOrBet_OnPreFlop(t *testing.T) {
	// In two player game, small blind acts first, he can't do Check and Bet
	table := table2Players_startedGame(t)
	sb := table.SmallBlindPosition()
	assert.NotNil(t, table.MakeAction(sb, getIden(sb), CheckAction))

	table = table2Players_startedGame(t)
	sb = table.SmallBlindPosition()
	assert.NotNil(t, table.MakeAction(sb, getIden(sb), BetAction(10)))
}

func TestAccept_SbCallOrRaiseOrFold_OnPreFlop(t *testing.T) {
	table := table2Players_startedGame(t)
	sb := table.SmallBlindPosition()
	assert.Nil(t, table.MakeAction(sb, getIden(sb), CallAction))

	table = table2Players_startedGame(t)
	sb = table.SmallBlindPosition()
	assert.Nil(t, table.MakeAction(sb, getIden(sb), FoldAction))
}


func TestReject_BbBetOrCall_AfterSbCall(t *testing.T) {
	table := table2Players_startedGame(t)
	table.tcall(t)
	bb := table.BigBlindPosition()
	assert.NotNil(t, table.MakeAction(bb, getIden(bb), CallAction))

	table = table2Players_startedGame(t)
	table.tcall(t)
	bb = table.BigBlindPosition()
	assert.NotNil(t, table.MakeAction(bb, getIden(bb), BetAction(20)))
}

func TestAccept_BbFoldOrCheckOrRaise_AfterSbCall(t *testing.T) {
	table := table2Players_startedGame(t)
	table.tcall(t)
	bb := table.BigBlindPosition()
	assert.Nil(t, table.MakeAction(bb, getIden(bb), FoldAction))

	table = table2Players_startedGame(t)
	table.tcall(t)
	bb = table.BigBlindPosition()
	assert.Nil(t, table.MakeAction(bb, getIden(bb), CheckAction))

	table = table2Players_startedGame(t)
	table.tcall(t)
	bb = table.BigBlindPosition()
	assert.Nil(t, table.MakeAction(bb, getIden(bb), RaiseAction(20)))
	assert.EqualValues(t, defaultBB + 20, table.GetBigBlind().TotalRoundBet)
}

func TestReject_BbCheckOrBet_AfterSbRaise(t *testing.T) {
	table := table2PlayersPositions_startedGame(t, 0, 1)
	err := table.MakeActionDeprecated(table.DecidingPosition, Raise, 20)
	assert.Nil(t, err)
	err = table.MakeActionDeprecated(table.DecidingPosition, Check, 0)
	assert.NotNil(t, err)

	table = table2PlayersPositions_startedGame(t, 0, 1)
	err = table.MakeActionDeprecated(table.DecidingPosition, Raise, 20)
	assert.Nil(t, err)
	err = table.MakeActionDeprecated(table.DecidingPosition, Bet, 0)
	assert.NotNil(t, err)
}

func TestAccept_BbFoldOrCallOrRaise_AfterSbRaise(t *testing.T) {
	table := table2PlayersPositions_startedGame(t, 0, 1)
	err := table.MakeActionDeprecated(table.DecidingPosition, Raise, 20)
	assert.Nil(t, err)
	err = table.MakeActionDeprecated(table.DecidingPosition, Fold, 0)
	assert.Nil(t, err)

	table = table2PlayersPositions_startedGame(t, 0, 1)
	err = table.MakeActionDeprecated(table.DecidingPosition, Raise, 20)
	assert.Nil(t, err)
	err = table.MakeActionDeprecated(table.DecidingPosition, Call, 0)
	assert.Nil(t, err)

	table = table2PlayersPositions_startedGame(t, 0, 1)
	err = table.MakeActionDeprecated(table.DecidingPosition, Raise, 20)
	assert.Nil(t, err)
	err = table.MakeActionDeprecated(table.DecidingPosition, Raise, 40)
	assert.Nil(t, err)
}


func TestReject_SbCallOrRaise_AfterBbCheck(t *testing.T) {
	table := table2Players_flop(t)
	err := table.MakeActionDeprecated(table.DecidingPosition, Call, 0)
	assert.NotNil(t, err)

	table = table2Players_flop(t)
	err = table.MakeActionDeprecated(table.DecidingPosition, Raise, 0)
	assert.NotNil(t, err)
}

func TestAccept_SbFoldOrCheck_Bet_AfterBbCheck(t *testing.T) {
	table := table2Players_flop(t)
	err := table.MakeActionDeprecated(table.DecidingPosition, Fold, 0)
	assert.Nil(t, err)

	table = table2Players_flop(t)
	err = table.MakeActionDeprecated(table.DecidingPosition, Check, 0)
	assert.Nil(t, err)

	table = table2Players_flop(t)
	err = table.MakeActionDeprecated(table.DecidingPosition, Bet, 10)
	assert.Nil(t, err)
}

func Test_RaceCondition_StaleTableOnTimeout(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table2Players_startedGame_chips(t, 0, 1, 250, 400)

	err := table.MakeAction(0, firstIdentity, RaiseAction(200))
	assert.Nil(t, err)

	err = table.MakeActionOnTimeout(0)
	assert.NotNil(t, err)
}

func TestRaiseZeroChips_ShouldDefaultToMin(t *testing.T) {
	table := table2PlayersPositions_startedGame(t, 0, 1)
	dp := table.DecidingPosition
	table.traise(t, 0)
	assert.EqualValues(t, defaultBB * 2, table.GetPlayerUnsafe(dp).TotalRoundBet)
}

func TestBetZeroChips_ShouldDefaultToMin(t *testing.T) {
	table := table2PlayersPositions_startedGame(t, 0, 1)
	table.tcall(t)
	table.tcheck(t)

	dp := table.DecidingPosition
	table.tbet(t, 0)
	assert.EqualValues(t, defaultBB, table.GetPlayerUnsafe(dp).TotalRoundBet)
}

func TestRaise_WhenBigBlindAllIn_WithLessChipsThanBB(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}
	table := table3Players_startedGame(t)
	table.tcall(t)
	table.tcall(t)
	table.tcheck(t)
	table.tfold(t)
	table.tbet(t, defaultBuyIn - defaultBB - 1)
	table.tcall(t)
	table.tcheck(t)
	table.tfold(t)
	assert.Nil(t, table.StartNextGame())

	table.tcall(t)
	assert.EqualValues(t, defaultBB, table.DealerPlayer().TotalRoundBet)

}
