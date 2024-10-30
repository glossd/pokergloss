package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_3P_AllIn_Acceptance(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}
	// Reject second player check, bet, and raise when not enough chips
	table := tableWithAllIn(t)
	err := table.MakeActionDeprecated(table.DecidingPosition, Check, 0)
	assert.NotNil(t, err)

	table = tableWithAllIn(t)
	err = table.MakeActionDeprecated(table.DecidingPosition, Bet, 10)
	assert.NotNil(t, err)

	table = tableWithAllIn(t)
	err = table.MakeActionDeprecated(table.DecidingPosition, Raise, 10)
	assert.NotNil(t, err)

	// Accept second player fold, call, allIn
	table = tableWithAllIn(t)
	err = table.MakeActionDeprecated(table.DecidingPosition, Fold, 0)
	assert.Nil(t, err)

	table = tableWithAllIn(t)
	err = table.MakeActionDeprecated(table.DecidingPosition, Call, 0)
	assert.Nil(t, err)


	table = tableWithAllIn(t)
	err = table.MakeActionDeprecated(table.DecidingPosition, AllIn, 0)
	assert.Nil(t, err)
}

func TestAllIn_RaiseAndBetTurnsIntoAllIn(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table2PlayersPositions_startedGame(t, 0, 1)
	err := table.MakeActionDeprecated(0, Raise, 249)
	assert.Nil(t, err)
	assert.EqualValues(t, AllIn, table.GetPlayerUnsafe(0).LastGameAction)

	table = table2PlayersPositions_startedGame(t, 0, 1)
	err = table.MakeActionDeprecated(0, Call, 0)
	assert.Nil(t, err)
	err = table.MakeActionDeprecated(1, Bet, 248)
	assert.Nil(t, err)
	assert.EqualValues(t, AllIn, table.GetPlayerUnsafe(1).LastGameAction)
}

func TestAllIn_DifferentStacks(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table3Players_startedGame_chips(t, 0, 1, 2, 250, 150, 250)
	err := table.MakeActionDeprecated(1, AllIn, 0)
	assert.Nil(t, err)
	err = table.MakeActionDeprecated(2, Raise, 200)
	assert.Nil(t, err)

	table = table3Players_startedGame_chips(t, 0, 1, 2, 250, 150, 250)
	err = table.MakeActionDeprecated(1, AllIn, 0)
	assert.Nil(t, err)
	err = table.MakeActionDeprecated(2, AllIn, 0)
	assert.Nil(t, err)
}

func tableWithAllIn(t *testing.T) *Table {
	table := table3Players_startedGame_chips(t, 0, 1, 2, 250, 150, 250)
	err := table.MakeActionDeprecated(table.DecidingPosition, AllIn, 0)
	assert.Nil(t, err)
	return table
}

func Test_3P_DealerFolds_ShouldNotKeepPlayingAfter(t *testing.T) {
	table := table3Players_startedGame(t)

	dealer := table.DecidingPlayerUnsafe()
	assert.EqualValues(t, Dealer, dealer.Blind)
	err := table.MakeActionDeprecated(dealer.Position, Fold, 0)
	assert.Nil(t, err)

	sb := table.DecidingPlayerUnsafe()
	assert.EqualValues(t, SmallBlind, sb.Blind)
	err = table.MakeActionDeprecated(sb.Position, Call, 0)
	assert.Nil(t, err)

	bb := table.DecidingPlayerUnsafe()
	assert.EqualValues(t, BigBlind, bb.Blind)
	err = table.MakeActionDeprecated(bb.Position, Check, 0)
	assert.Nil(t, err)

	assert.True(t, table.IsFlop())

	assert.EqualValues(t, sb.Position, table.DecidingPosition)
	err = table.MakeActionDeprecated(sb.Position, Check, 0)
	assert.Nil(t, err)
	err = table.MakeActionDeprecated(bb.Position, Check, 0)
	assert.Nil(t, err)

	assert.True(t, table.IsTurn())
	assert.EqualValues(t, sb.Position, table.DecidingPosition)
}

func Test_3P_FoldedPlayer_ShouldNotKeepPlaying_AfterBet(t *testing.T) {
	table := table3Players_startedGame(t)

	foldedPosition := table.DecidingPosition
	table.tfold(t)
	table.tcall(t)
	table.traise(t, 10)
	assert.NotEqualValues(t, foldedPosition, table.DecidingPosition)
}

func Test_3P_AllTimeout(t *testing.T) {
	table := table3Players_startedGame(t)

	err := table.MakeActionOnTimeout(table.DecidingPosition)
	assert.Nil(t, err)
	err = table.MakeActionOnTimeout(table.DecidingPosition)
	assert.Nil(t, err)

	assert.True(t, table.IsGameEnd())
	assert.Len(t, table.Winners, 1)
}

func Test_2P_SbAllIn(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table2Players_startedGame_chips(t, 0, 1, 250, 400)
	table.tallIn(t)
	table.tcall(t)

	assert.True(t, table.IsGameEnd())
	assert.Len(t, table.CommunityCards.AvailableCards(), 5)
}

func Test_2P_SbAllIn_BbFolds(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table2Players_startedGame_chips(t, 0, 1, 250, 400)
	table.tallIn(t)
	table.tfold(t)

	assert.True(t, table.IsGameEnd())
	assert.Len(t, table.CommunityCards.AvailableCards(), 0)
}

func Test_3P_SecondPlayerAllIn_ShouldSkipHim(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}
	table := table3Players_startedGame_chips(t, 0, 1, 2, 250, 150, 250)

	err := table.MakeActionDeprecated(1, AllIn, 0)
	assert.Nil(t, err)

	err = table.MakeActionDeprecated(2, Call, 0)
	assert.Nil(t, err)

	err = table.MakeActionDeprecated(0, Call, 0)
	assert.Nil(t, err)

	// Flop
	assert.True(t, table.IsNewRound())

	checkRound3P(t, table)

	// Turn
	assert.True(t, table.IsNewRound())

	checkRound3P(t, table)

	// River
	assert.True(t, table.IsNewRound())

	checkRound3P(t, table)

	assert.True(t, table.IsGameEnd())
	assert.EqualValues(t, PlayerPlaying, table.GetPlayerUnsafe(1).Status)

	assert.Nil(t, table.StartNextGame())
}

func Test_3P_StaysWithChips_AllInWin_AllInLoose(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = AlgoMock_3POrdered_Second_Third_First(t)

	table := table3Players_startedGameOrder_chips(t, 0, 1, 2, 177, 227, 396)
	assert.EqualValues(t, 2, table.BigBlindPosition())
	assert.EqualValues(t, 0, table.DealerPosition())
	table.traise(t, 5)
	table.tcall(t)
	table.tcall(t)                     //Stacks
	assert.True(t, table.IsNewRound()) //6, 6 ,6
	table.tbet(t, 6) //   0, 6, 0
	table.traise(t, 8) // 0, 6, 8
	table.traise(t, 12) //12,6,8
	table.traise(t, 7)  //12,13,8
	table.traise(t, 13) //12,13,21
	table.tallIn(t) // first player all-in
	table.traise(t, 190)//all,203,21
	table.traise(t, 300)//all,203,321

	table.tallIn(t) // second player all-in

	assert.True(t, table.IsGameEnd())
	assert.Len(t, table.Pots, 2)

	assert.EqualValues(t, 3*177, table.Pots[0].Chips)
	assert.Len(t, table.Pots[0].UserIDs, 1)
	assert.EqualValues(t, firstIdentity.UserId, table.Pots[0].UserIDs[0])

	assert.EqualValues(t, 2*(227 - 177), table.Pots[1].Chips)
	assert.Len(t, table.Pots[1].UserIDs, 1)
	assert.EqualValues(t, secondIdentity.UserId, table.Pots[1].UserIDs[0])

	assert.Len(t, table.Winners, 1)
	assert.EqualValues(t, 3*177 + 2*(227-177), table.Winners[0].Chips)
	assert.EqualValues(t, 1, table.Winners[0].Position)
}

func checkRound3P(t *testing.T, table *Table) {
	err := table.MakeActionDeprecated(2, Check, 0)
	assert.Nil(t, err)

	err = table.MakeActionDeprecated(0, Check, 0)
	assert.Nil(t, err)
}

func Test_2P_PlayerWithLessChipsAllIn(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)

	fPos := 0
	sPos := 1

	table := table2P_firstAllIn(t)
	assert.NotNil(t, table.MakeActionDeprecated(sPos, Check, 0))
	table = table2P_firstAllIn(t)
	assert.NotNil(t, table.MakeActionDeprecated(sPos, Raise, 150))
	table = table2P_firstAllIn(t)
	assert.NotNil(t, table.MakeActionDeprecated(sPos, Bet, 150))
	table = table2P_firstAllIn(t)
	assert.NotNil(t, table.MakeActionDeprecated(sPos, AllIn, 0))

	table = table2P_firstAllIn(t)
	assert.Nil(t, table.MakeActionDeprecated(sPos, Fold, 0))
	assert.EqualValues(t, 102, table.GetPlayerUnsafe(fPos).Stack)
	assert.EqualValues(t, 248, table.GetPlayerUnsafe(sPos).Stack)

	table = table2P_firstAllIn(t)
	assert.Nil(t, table.MakeActionDeprecated(sPos, Call, 0))

	assert.EqualValues(t, 200, table.GetPlayerUnsafe(fPos).Stack)
	assert.EqualValues(t, 150, table.GetPlayerUnsafe(sPos).Stack)
}

func Test_2P_PlayerGetSomeBackAfterSmallAllIn(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = Algo_2P_MockSecondPlayerLoses()
	table := table2Players_startedGame_chips(t, 0, 1, 400, 100)
	table.tallIn(t)
	table.tallIn(t)
	assert.EqualValues(t, 200, table.TotalPot)
	assert.EqualValues(t, 1, len(table.Pots))
	assert.EqualValues(t, 200, table.Pots[0].Chips)
}

func TestNewCommunityCardsOnPreFlop(t *testing.T) {
	table := table2Players_startedGame(t)
	table.tallIn(t)
	table.tallIn(t)
	assert.EqualValues(t, 5, len(table.CommunityCards.GetNewCards()))
}

func table2P_firstAllIn(t *testing.T) *Table {
	Algo = Algo_2P_MockSecondPlayerLoses()

	table := table2Players_startedGame_chips(t, 0, 1, 100, 250)
	err := table.MakeActionDeprecated(0, AllIn, 0)
	assert.Nil(t, err)
	return table
}
