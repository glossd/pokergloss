package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_2P_Winners(t *testing.T) {
	table := table2PlayerPositions_gameRiver(t, 0, 1)
	sbPosition := table.DecidingPosition
	err := table.MakeActionDeprecated(sbPosition, Check, 0)
	assert.Nil(t, err)
	bbPosition := table.DecidingPosition
	err = table.MakeActionDeprecated(bbPosition, Check, 0)
	assert.Nil(t, err)

	assert.Contains(t, []int{1, 2}, len(table.Winners), "Only 1 or two winners allowed in 2 player game")
	if len(table.Winners) == 1 {
		winnerPosition := table.Winners[0].Position
		assert.EqualValues(t, 252, table.GetPlayerUnsafe(winnerPosition).Stack)
		loserPosition := notPlayerPosition(table.GetPlayerUnsafe(winnerPosition), 0, 1)
		assert.EqualValues(t, 248, table.GetSeatUnsafe(loserPosition).GetPlayer().Stack)
	}
	if len(table.Winners) == 2 {
		assert.EqualValues(t, 250, table.GetPlayerUnsafe(sbPosition).Stack)
		assert.EqualValues(t, 250, table.GetPlayerUnsafe(bbPosition).Stack)
	}
}

func Test_2P_WinnersOnTimeout(t *testing.T) {
	table := table2PlayerPositions_gameRiver(t, 0, 1)
	sbPosition := table.DecidingPosition
	err := table.MakeActionDeprecated(sbPosition, Check, 0)
	assert.Nil(t, err)
	bbPosition := table.DecidingPosition
	err = table.MakeActionOnTimeout(bbPosition)
	assert.Nil(t, err)

	assert.EqualValues(t, 1, len(table.Winners))
}

func TestTable_ComputeWinners(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = Algo_2P_MockSecondPlayerLoses()
	table := table2Players_startedGame_chips(t, 0, 1, 250, 250)
	table.traise(t, 9)
	table.tfold(t)

	assert.Len(t, table.Winners, 1)
	assert.Len(t, table.Pots, 1)
	assert.Len(t, table.Pots[0].WinnerPositions, 1)
	assert.EqualValues(t, 0, table.Pots[0].WinnerPositions[0])
	assert.EqualValues(t, table.GetPlayerUnsafe(0).Stack, 252)
	assert.EqualValues(t, table.GetPlayerUnsafe(1).Stack, 248)
}

func Test_3P_TwoAllInWithSameStacksSameRank(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = AlgoMock_3P_ThirdAndFirst_Second(t)
	table := table3Players_startedGame_chips(t, 0, 1, 2, 250, 350, 250)
	assert.EqualValues(t, 1, table.DecidingPosition)
	table.tcall(t)
	table.tallIn(t)
	table.tallIn(t)
	table.tcall(t)

	assert.Len(t, table.Winners, 2)
	assert.Len(t, table.Pots, 1)
	assert.EqualValues(t, table.Pots[0].Chips, 750)
	assert.Len(t, table.Pots[0].WinnerPositions, 2)
	assert.EqualValues(t, table.Pots[0].WinnerPositions[0], 0)
	assert.EqualValues(t, table.Pots[0].WinnerPositions[1], 2)

	assert.EqualValues(t, 375, table.GetPlayerUnsafe(0).Stack)
	assert.EqualValues(t, 100, table.GetPlayerUnsafe(1).Stack)
	assert.EqualValues(t, 375, table.GetPlayerUnsafe(2).Stack)
}

func Test_3P_TwoAllInDiffStacksDiffRank(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = AlgoMock_3P_Second_Third_First(t)

	table := table3Players_startedGame_chips(t, 0, 1, 2, 350, 150, 250)
	assert.EqualValues(t, 1, table.DecidingPosition)
	table.tallIn(t)
	table.tcall(t)
	table.tcall(t)
	assert.True(t, table.IsNewRound())

	table.tallIn(t)
	table.tcall(t)
	assert.True(t, table.IsGameEnd())

	assert.Len(t, table.Winners, 2)
	assert.EqualValues(t, 1, table.Winners[0].Position)
	assert.EqualValues(t, 450, table.Winners[0].Chips)
	assert.EqualValues(t, 2, table.Winners[1].Position)
	assert.EqualValues(t, 200, table.Winners[1].Chips)

	assert.Len(t, table.Pots, 2)
	assert.Len(t, table.Pots[0].WinnerPositions, 1)
	assert.EqualValues(t, 1, table.Pots[0].WinnerPositions[0])
	assert.Len(t, table.Pots[1].WinnerPositions, 1)
	assert.EqualValues(t, 2, table.Pots[1].WinnerPositions[0])

	assert.EqualValues(t, 450, table.GetPlayerUnsafe(1).Stack)
	assert.EqualValues(t, 200, table.GetPlayerUnsafe(2).Stack)
	assert.EqualValues(t, 100, table.GetPlayerUnsafe(0).Stack)
}

func Test_3P_OneAllInOneFoldDiffRank(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = AlgoMock_3POrdered_Second_Third_First(t)

	table := table3Players_startedGameOrder_chips(t, 0, 1, 2, 250, 150, 250)
	assert.EqualValues(t, 0, table.DecidingPosition)
	table.tcall(t)
	table.tcall(t)
	table.tfold(t)

	table.tallIn(t)
	table.tcall(t)

	assert.Len(t, table.Winners, 1)
	assert.EqualValues(t, 1, table.Winners[0].Position)
	assert.EqualValues(t, 302, table.Winners[0].Chips)

	assert.Len(t, table.Pots, 1)

	assert.EqualValues(t, 100, table.GetPlayerUnsafe(0).Stack)
	assert.EqualValues(t, 302, table.GetPlayerUnsafe(1).Stack)
	assert.EqualValues(t, 248, table.GetPlayerUnsafe(2).Stack)
}

func Test_3P_FoldBiggerThanAllIn(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = AlgoMock_3POrdered_Second_Third_First(t)

	table := table3Players_startedGameOrder_chips(t, 0, 1, 2, 350, 150, 250)
	table.traise(t, 200)
	table.tallIn(t)
	table.tallIn(t)
	table.tfold(t)

	assert.Len(t, table.Winners, 2)
	assert.EqualValues(t, 1, table.Winners[0].Position)
	assert.EqualValues(t, 450, table.Winners[0].Chips)
	assert.EqualValues(t, 2, table.Winners[1].Position)
	assert.EqualValues(t, 100, table.Winners[1].Chips)

	assert.Len(t, table.Pots, 2) // i'm confused, there should be third pot with zero chips
	assert.EqualValues(t, 450, table.Pots[0].Chips)
	assert.EqualValues(t, 100, table.Pots[1].Chips)


	assert.EqualValues(t, 150, table.GetPlayerUnsafe(0).Stack)
	assert.EqualValues(t, 450, table.GetPlayerUnsafe(1).Stack)
	assert.EqualValues(t, 150, table.GetPlayerUnsafe(2).Stack)
}

func Test_4P_FirstTwoAllInsWithSameChipsSameRankThenAnotherAllIn(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = AlgoMock_4P_FirstAndSecond_Third_Fourth(t)
	table := table4PlayersPositions_startedGame_chips(t, 0, 1, 2, 3, 150, 150, 200, 250)

	assert.Nil(t, table.MakeActionDeprecated(0, AllIn, 0))
	assert.Nil(t, table.MakeActionDeprecated(1, AllIn, 0))
	assert.Nil(t, table.MakeActionDeprecated(2, Call, 0))
	assert.Nil(t, table.MakeActionDeprecated(3, Call, 0))

	assert.Nil(t, table.MakeActionDeprecated(2, AllIn, 0))
	assert.Nil(t, table.MakeActionDeprecated(3, Call, 0))

	assert.Len(t, table.Winners, 3)
	assert.EqualValues(t, 300, table.Winners[0].Chips)
	assert.EqualValues(t, 300, table.Winners[1].Chips)
	assert.EqualValues(t, 100, table.Winners[2].Chips)

	assert.Len(t, table.Pots, 2)

	assert.EqualValues(t, 300, table.GetPlayerUnsafe(0).Stack)
	assert.EqualValues(t, 300, table.GetPlayerUnsafe(1).Stack)
	assert.EqualValues(t, 100, table.GetPlayerUnsafe(2).Stack)
	assert.EqualValues(t, 50, table.GetPlayerUnsafe(3).Stack)
}
