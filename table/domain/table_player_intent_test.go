package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntent_OnNewRound(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table2Players_flop(t)
	assert.Nil(t, table.SetIntent(0, firstIdentity, CheckIntent))
	assert.Nil(t, table.MakeAction(1, secondIdentity, CheckAction))

	assert.True(t, table.IsNewRound())
	assert.EqualValues(t, table.DecidingPosition, 1)
	assert.Nil(t, table.GetPlayerUnsafe(0).Intent)

	assert.Len(t, table.MadeActionPlayerPositions(), 2)
}

func TestIntent_PreFlopBigBlind(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table2PlayersPositions_startedGame(t, 0, 1)
	assert.Nil(t, table.SetIntent(1, secondIdentity, CheckIntent))
	assert.Nil(t, table.MakeActionDeprecated(0, Call, 0))

	assert.True(t, table.IsNewRound())
	assert.EqualValues(t, table.DecidingPosition, 1)
	assert.Len(t, table.MadeActionPlayerPositions(), 2)
	assert.Nil(t, table.GetPlayerUnsafe(1).Intent)
}

func TestIntent_Raise(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table2Players_startedGame(t)
	bb := table.BigBlindPosition()
	assert.Nil(t, table.SetIntent(bb, getIden(bb), NewIntent(RaiseIntentType, defaultBB)))
	table.tcall(t)

	assert.Len(t, table.MadeActionPlayerPositions(), 2)
	assert.Nil(t, table.GetPlayerUnsafe(1).Intent)
}

func TestIntent_OnBet(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table3Players_flop(t)
	intent := NewIntent(CheckIntentType, 0)
	assert.Nil(t, table.SetIntent(0, firstIdentity, intent))
	assert.Nil(t, table.MakeActionDeprecated(1, Check, 0))
	assert.Nil(t, table.MakeActionDeprecated(2, Check, 0))

	assert.True(t, table.IsNewRound())
	assert.EqualValues(t, table.DecidingPosition, 1)
	assert.Len(t, table.MadeActionPlayerPositions(), 2)
	assert.Nil(t, table.GetPlayerUnsafe(0).Intent)
}

func TestIntent_Multiple(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table3Players_flop(t)
	assert.Nil(t, table.SetIntent(0, firstIdentity, CheckIntent))
	assert.Nil(t, table.SetIntent(2, thirdIdentity, CheckIntent))
	assert.Nil(t, table.MakeAction(1, secondIdentity, CheckAction))

	assert.True(t, table.IsNewRound())
	assert.EqualValues(t, table.DecidingPosition, 1)
	assert.Len(t, table.MadeActionPlayerPositions(), 3)
}

func TestIntent_UpgradeOnNewRound(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table2Players_flop(t)
	intent := NewIntent(CheckCallAnyIntentType, 0)
	assert.Nil(t, table.SetIntent(0, firstIdentity, intent))
	assert.Nil(t, table.MakeActionDeprecated(1, Bet, 5))

	assert.True(t, table.IsNewRound())
	assert.EqualValues(t, table.DecidingPosition, 1)
	assert.Len(t, table.MadeActionPlayerPositions(), 2)
	assert.EqualValues(t, table.GetPlayerUnsafe(0).LastGameAction, Call)
}

func TestIntent_RemoveIntentAfterBet(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table2Players_flop(t)
	assert.Nil(t, table.SetIntent(0, firstIdentity, CheckIntent))
	assert.Nil(t, table.MakeAction(1, secondIdentity, BetAction(10)))

	assert.False(t, table.IsNewRound())
	assert.EqualValues(t, table.DecidingPosition, 0)
}

func TestIntent_BusinessAcceptReject(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}
	table := table4Players_startedGame(t)
	assert.Nil(t, table.SetIntent(1, secondIdentity, CallIntent))
	assert.Nil(t, table.SetIntent(2, thirdIdentity, CallIntent))
	assert.Nil(t, table.SetIntent(3, fourthIdentity, CheckIntent))
	assert.Nil(t, table.MakeAction(0, firstIdentity, CallAction))

	assert.True(t, table.IsNewRound())

	assert.EqualValues(t, SmallBlind, table.DecidingPlayerUnsafe().Blind)
	assert.EqualValues(t, 2, table.DecidingPosition)

	assert.NotNil(t, table.SetIntent(2, thirdIdentity, CheckIntent))
	assert.Nil(t, table.SetIntent(3, fourthIdentity, CheckIntent))
	assert.Nil(t, table.RemoveIntent(3, fourthIdentity))
	assert.Nil(t, table.SetIntent(0, firstIdentity, CheckIntent))
	assert.Nil(t, table.RemoveIntent(0, firstIdentity))
	assert.Nil(t, table.SetIntent(1, secondIdentity, CheckIntent))
	assert.Nil(t, table.RemoveIntent(1, secondIdentity))

	assert.Nil(t, table.MakeAction(2, thirdIdentity, CheckAction))
	assert.NotNil(t, table.SetIntent(2, thirdIdentity, CheckIntent))
	assert.Nil(t, table.SetIntent(0, firstIdentity, CheckIntent))
	assert.Nil(t, table.SetIntent(1, secondIdentity, CheckIntent))
}

func TestIntent_BusinessBothDirections(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}
	table := table2PlayersPositions_startedGame(t, 0, 8)
	assert.Nil(t, table.SetIntent(8, secondIdentity, CheckIntent))
	assert.Nil(t, table.RemoveIntent(8, secondIdentity))
	assert.Nil(t, table.MakeAction(0, firstIdentity, CallAction))
	assert.NotNil(t, table.SetIntent(0, firstIdentity, CheckIntent))

	assert.Nil(t, table.MakeAction(8, secondIdentity, CheckAction))
	assert.Nil(t, table.SetIntent(0, firstIdentity, CheckIntent))

	assert.Nil(t, table.MakeAction(8, secondIdentity, FoldAction))

	// new game
	assert.Nil(t, table.StartNextGame())

	assert.Nil(t, table.SetIntent(0, firstIdentity, CheckIntent))
	assert.Nil(t, table.RemoveIntent(0, firstIdentity))
	assert.Nil(t, table.MakeAction(8, secondIdentity, CallAction))
	assert.NotNil(t, table.SetIntent(8, secondIdentity, CheckIntent))

	assert.Nil(t, table.MakeAction(0, firstIdentity, CheckAction))
	assert.Nil(t, table.SetIntent(8, secondIdentity, CheckIntent))
	assert.Nil(t, table.RemoveIntent(8, secondIdentity))
	assert.Nil(t, table.MakeAction(0, firstIdentity, CheckAction))
	assert.NotNil(t, table.SetIntent(0, firstIdentity, CheckIntent))
}

func TestIntent_2P_AllInIntentOnRiver(t *testing.T) {
	table := table2Players_startedGame_chips(t, 0, 1, 150, 300)
	table.tcall(t)
	table.tcheck(t)
	assert.True(t, table.IsNewRound())
	table.tcheck(t)
	table.tcheck(t)
	assert.True(t, table.IsNewRound())
	table.tcheck(t)
	table.tbet(t, 20)
	table.tcall(t)
	assert.True(t, table.IsNewRound())
	table.tsetIntent(t, table.nextDecidingPos_2P(), AllInIntent)
	table.tcheck(t)
	table.tfold(t)
}

func TestIntent_AcceptReject_NewRound(t *testing.T) {
	acceptIntentTableOnNewRound(t, FoldIntentType, 0)
	acceptIntentTableOnNewRound(t, AllInIntentType, 0)

	acceptIntentTableOnNewRound(t, CheckFoldIntentType, 0)
	acceptIntentTableOnNewRound(t, CheckIntentType, 0)
	acceptIntentTableOnNewRound(t, CheckCallAnyIntentType, 0)
	acceptIntentTableOnNewRound(t, BetIntentType, 10)

	rejectIntentTableOnNewRound(t, CallIntentType, 0)
	rejectIntentTableOnNewRound(t, CallFoldIntentType, 0)
	rejectIntentTableOnNewRound(t, CallAnyIntentType, 0)
	rejectIntentTableOnNewRound(t, RaiseIntentType, 10)
}

func acceptIntentTableOnNewRound(t *testing.T, intentType IntentType, chips int64) {
	intent := NewIntent(intentType, chips)
	table := table2Players_flop(t)
	nextP := table.nextPlayingPlayerUnsafe(table.DecidingPosition)
	assert.Nil(t, table.SetIntent(nextP.Position, nextP.Identity, intent))
}

func rejectIntentTableOnNewRound(t *testing.T, intentType IntentType, chips int64) {
	intent := NewIntent(intentType, chips)
	table := table2Players_flop(t)
	nextP := table.nextPlayingPlayerUnsafe(table.DecidingPosition)
	assert.NotNil(t, table.SetIntent(nextP.Position, nextP.Identity, intent))
}


func TestIntent_AcceptReject_OnBet(t *testing.T) {
	acceptIntentTableOnBet(t, FoldIntentType, 0)
	acceptIntentTableOnBet(t, AllInIntentType, 0)

	acceptIntentTableOnBet(t, CallIntentType, 0)
	acceptIntentTableOnBet(t, CallFoldIntentType, 0)
	acceptIntentTableOnBet(t, CallAnyIntentType, 0)
	acceptIntentTableOnBet(t, RaiseIntentType, 10)

	rejectIntentTableOnBet(t, CheckFoldIntentType, 0)
	rejectIntentTableOnBet(t, CheckIntentType, 0)
	rejectIntentTableOnBet(t, CheckCallAnyIntentType, 0)
	rejectIntentTableOnBet(t, BetIntentType, 10)
}

func acceptIntentTableOnBet(t *testing.T, intentType IntentType, chips int64) {
	intent := NewIntent(intentType, chips)
	table := table3Players_startedGame(t)
	nextP := table.nextPlayingPlayerUnsafe(table.DecidingPosition)
	assert.Nil(t, table.SetIntent(nextP.Position, nextP.Identity, intent))
}

func rejectIntentTableOnBet(t *testing.T, intentType IntentType, chips int64) {
	intent := NewIntent(intentType, chips)
	table := table3Players_startedGame(t)
	nextP := table.nextPlayingPlayerUnsafe(table.DecidingPosition)
	assert.NotNil(t, table.SetIntent(nextP.Position, nextP.Identity, intent))
}

func TestBetIntent_ShouldConvertToAllIn(t *testing.T) {
	table := table2Players_startedGame(t)

	intentPos := table.DecidingPosition
	table.tcall(t)
	table.tcheck(t)

	err := table.SetIntent(intentPos, getIden(intentPos), NewIntent(BetIntentType, defaultBuyIn-table.BigBlind))
	assert.Nil(t, err)
	table.tcheck(t)
	assert.EqualValues(t, AllIn, table.GetPlayerUnsafe(intentPos).LastGameAction)
}

func TestRaiseIntent_ShouldConvertToAllIn(t *testing.T) {
	table := table3Players_startedGame(t)
	sbPos := table.SmallBlindPosition()
	err := table.SetIntent(sbPos, getIden(sbPos), NewIntent(RaiseIntentType, table.GetPlayerUnsafe(sbPos).Stack))
	assert.Nil(t, err)
	table.tcall(t)
	assert.EqualValues(t, AllIn, table.GetPlayerUnsafe(sbPos).LastGameAction)
}

func TestRaiseIntent_CanNotBeMoreThanStack(t *testing.T) {
	Algo = &MockAlgo{}
	table := table3Players_startedGame_chips(t, 0, 1, 2, 100, 400, defaultBuyIn)
	assert.EqualValues(t, 1, table.DealerPosition())
	table.tallIn(t)
	assert.NotNil(t, table.SetIntent(0, getIden(0), NewIntent(RaiseIntentType, 800)))
}

func TestBetIntent_CanNotBeMoreThanStack(t *testing.T) {
	Algo = &MockAlgo{}
	table := table3Players_startedGame(t)
	assert.EqualValues(t, 1, table.DealerPosition())
	table.tcall(t)
	table.tcall(t)
	table.tcheck(t)
	assert.NotNil(t, table.SetIntent(1, getIden(1), NewIntent(BetIntentType, 1000)))
}

func TestCallIntent_CanNotBeMoreThanStack(t *testing.T) {
	Algo = &MockAlgo{dealerPos: 1}
	table := table3Players_startedGameOrder_chips(t, 0, 1, 2, 400, 250, 100)
	assert.EqualValues(t, 0, table.DealerPosition())
	table.tallIn(t)
	err := table.SetIntent(2, getIden(2), CallIntent)
	assert.NotNil(t, err)
	table.tfold(t)
}
