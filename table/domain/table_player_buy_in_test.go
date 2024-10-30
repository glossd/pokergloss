package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_2P_GameStarted(t *testing.T) {
	table := defaultTable(t)
	err := table.ReserveSeat(0, firstIdentity)
	assert.Nil(t, err)
	isNewGame, err := table.BuyIn(250, 0, firstIdentity)
	assert.Nil(t, err)
	assert.False(t, isNewGame)

	err = table.ReserveSeat(1, secondIdentity)
	assert.Nil(t, err)
	newGame, err := table.BuyIn(250, 1, secondIdentity)
	assert.Nil(t, err)
	assert.True(t, newGame)
}

func TestRejectPlayerPuttingChipsOnAnotherPosition(t *testing.T) {
	table := defaultTable(t)
	err := table.ReserveSeat(0, firstIdentity)
	assert.Nil(t, err)
	isNewGame, err := table.BuyIn(250, 7, firstIdentity)
	assert.False(t, isNewGame)
	assert.NotNil(t, err)
	assert.Equal(t, "table position is not yours", err.Error())
}

func TestGameShouldNotStart_OnThirdPlayerJoining(t *testing.T) {
	table := table2Players_startedGame(t)
	err := table.ReserveSeat(2, getIden(2))
	assert.Nil(t, err)
	isNewGame, err := table.BuyIn(250, 2, getIden(2))
	assert.Nil(t, err)
	assert.False(t, isNewGame)
}

func TestBuyIn_ShouldFail_WhileReady(t *testing.T) {
	table := defaultTable(t)
	err := table.ReserveSeat(0, getIden(0))
	assert.Nil(t, err)
	_, err = table.BuyIn(defaultBuyIn, 0, getIden(0))
	assert.Nil(t, err)

	_, err = table.BuyIn(100, 0, getIden(0))
	assert.NotNil(t, err)
}

func TestBuyIn_ShouldFail_WhilePlaying(t *testing.T) {
	table := table2Players_startedGame(t)
	_, err := table.BuyIn(defaultBuyIn, 0, getIden(0))
	assert.NotNil(t, err)
}

func TestBuyIn_ShouldFail_WhileSittingOutWithStack(t *testing.T) {
	table := table2Players_startedGame(t)
	timeoutPos := table.DecidingPosition
	assert.Nil(t, table.MakeActionOnTimeout(timeoutPos))
	assert.EqualValues(t, PlayerSittingOut, table.GetPlayerUnsafe(timeoutPos).Status)

	_, err := table.BuyIn(defaultBuyIn, timeoutPos, getIden(timeoutPos))
	assert.NotNil(t, err)
}

func TestBuyIn_ShouldSucceedWhenBroke(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = Algo_2P_MockFirstPlayerLoses()
	table := table2Players_startedGame(t)
	table.tallIn(t)
	table.tallIn(t)
	assert.Nil(t, table.StartNextGame())
	assert.EqualValues(t, PlayerReservedSeat, table.GetPlayerUnsafe(0).Status)

	isNewGame, err := table.BuyIn(defaultBuyIn, 0, getIden(0))
	assert.True(t, isNewGame)
	assert.Nil(t, err)
}

func TestBuyIn_ShouldCancelWhenBroke(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = Algo_2P_MockFirstPlayerLoses()
	table := table2Players_startedGame(t)
	table.tallIn(t)
	table.tallIn(t)
	assert.Nil(t, table.StartNextGame())
	assert.EqualValues(t, PlayerReservedSeat, table.GetPlayerUnsafe(0).Status)

	assert.Nil(t, table.CancelSeatReservation(0, getIden(0)))
}
