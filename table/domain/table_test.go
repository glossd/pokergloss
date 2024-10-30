package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecidingPlayerPosition(t *testing.T) {
	table := table2PlayersPositions_startedGame(t, 0, 1)
	sbPosition := table.DecidingPosition
	assert.Contains(t, []int{0, 1}, sbPosition)

	err := table.MakeActionDeprecated(sbPosition, Call, 0)
	assert.Nil(t, err)
	bbPosition := table.DecidingPosition
	assert.Contains(t, []int{0, 1}, sbPosition)
	assert.NotEqualValues(t, sbPosition, bbPosition)
}

func TestTablePot(t *testing.T) {
	table := table2Players_startedGame(t)
	table.tcall(t)
	table.tcheck(t)
	assert.EqualValues(t, 1, len(table.Pots))
	assert.EqualValues(t, 4, table.Pots[0].Chips)
	assert.EqualValues(t, 0, len(table.Pots[0].WinnerPositions))
}

func TestCheckNewRound(t *testing.T) {
	table, sb := table2Players_startedGameDeprecated(t, 0, 1)
	err := table.MakeActionDeprecated(sb.Position, Call, 0)
	assert.Nil(t, err)
	err = table.MakeActionDeprecated(table.DecidingPosition, Check, 0)
	assert.Nil(t, err)
	assert.True(t, table.IsNewRound())
	assert.True(t, table.IsFlop())
}

func Test_DecidingPosition_StoppedTable(t *testing.T) {
	fPosition := 0
	sPosition := 1
	table, sb := table2Players_startedGameDeprecated(t, fPosition, sPosition)

	err := table.MakeActionOnTimeout(sb.Position)
	assert.Nil(t, err)

	assert.Nil(t, table.StartNextGame())

	assert.EqualValues(t, -1, table.DecidingPosition)
}

func TestTable_IsWaitingAfterStopTable(t *testing.T) {
	fPosition := 0
	sPosition := 1
	table, sb := table2Players_startedGameDeprecated(t, fPosition, sPosition)

	err := table.MakeActionOnTimeout(sb.Position)
	assert.Nil(t, err)

	assert.Nil(t, table.StartNextGame())

	assert.True(t, table.IsWaiting())
}

func Test_CommunityCards_StopTable(t *testing.T) {
	fPosition := 0
	sPosition := 1
	table := table2PlayerPositions_gameRiver(t, fPosition, sPosition)

	err := table.MakeActionOnTimeout(table.DecidingPosition)
	assert.Nil(t, err)

	assert.Nil(t, table.StartNextGame())

	assert.Nil(t, table.CommunityCards.Flop)
	assert.Nil(t, table.CommunityCards.Turn)
	assert.Nil(t, table.CommunityCards.River)
}

func TestPotOnBet(t *testing.T) {
	table, _ := table2Players_startedGameDeprecated(t, 0, 1)
	assert.EqualValues(t, 1, len(table.Pots))
	assert.Zero(t, table.Pots[0].Chips)
	assert.EqualValues(t, 3, table.TotalPot)

	table.traise(t, 3)
	assert.Zero(t, table.Pots[0].Chips)
	assert.EqualValues(t, 6, table.TotalPot)

	table.tcall(t)
	assert.EqualValues(t, 8, table.Pots[0].Chips)
	assert.EqualValues(t, 8, table.TotalPot)
}

func TestPotAtTheEnd(t *testing.T) {
	table := table2PlayerPositions_gameRiver(t, 0, 1)
	table.tbet(t, 4)
	table.tcall(t)
	assert.True(t, table.IsGameEnd())
	assert.EqualValues(t, 12, table.Pots[0].Chips)
	assert.EqualValues(t, 1, len(table.Pots))
}

func Test_2P_GameEnded(t *testing.T) {
	table := table2PlayerPositions_gameRiver(t, 0, 1)
	sbPosition := table.DecidingPosition
	bbPosition := notPlayerPosition(table.DecidingPlayerUnsafe(), 0, 1)
	err := table.MakeActionDeprecated(sbPosition, Check, 0)
	assert.Nil(t, err)
	err = table.MakeActionDeprecated(bbPosition, Check, 0)
	assert.Nil(t, err)

	assert.True(t, table.IsGameEnd())
}

func Test_2P_StoppedTable(t *testing.T) {
	table, sb := table2Players_startedGameDeprecated(t, 0, 1)
	err := table.MakeActionOnTimeout(sb.Position)
	assert.Nil(t, err)
	assert.True(t, table.IsGameEnd())
}

func TestTable_GetStackOverflowPlayer(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = Algo_2P_MockSecondPlayerLoses()
	table := table2Players_startedGame_chips(t, 0, 1, 250, 100)
	table.tallIn(t)
	table.tallIn(t)
	sop := table.BuildStackOverflowPlayer()
	assert.EqualValues(t, 150, sop.Stack)

	assert.Nil(t, table.StartNextGame())
	assert.EqualValues(t, 350, table.GetPlayerUnsafe(0).Stack)
}

func notPlayerPosition(decidingP *Player, fPosition int, sPosition int) int {
	if decidingP.Position == fPosition {
		return  sPosition
	} else {
		return fPosition
	}
}
