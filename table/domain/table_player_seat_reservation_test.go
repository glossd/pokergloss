package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRejectAddingPlayerOnTakenSeat(t *testing.T) {
	table := defaultTable(t)
	err := table.ReserveSeat(0, firstIdentity)
	assert.Nil(t, err)
	err = table.ReserveSeat(0, secondIdentity)
	assert.NotNil(t, err)
	assert.Equal(t, "the seat is taken", err.Error())
}

func TestRejectPlayerTakingAnotherSeat(t *testing.T) {
	table := defaultTable(t)
	err := table.ReserveSeat(0, firstIdentity)
	assert.Nil(t, err)
	err = table.ReserveSeat(1, firstIdentity)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAlreadySitting, err)
}

func TestRejectPlayerCancelReservationWhilPlayer(t *testing.T) {
	table := table2Players_startedGame(t)
	assert.NotNil(t, table.CancelSeatReservation(0, firstIdentity))
	assert.NotNil(t, table.CancelSeatReservationTimeout(0))
}

func TestReserveAndCancel(t *testing.T) {
	table := defaultTable(t)

	assert.Nil(t, table.ReserveSeat(0, firstIdentity))
	assert.NotNil(t, table.GetPlayerUnsafe(0))

	assert.Nil(t, table.CancelSeatReservation(0, firstIdentity))
	assert.Nil(t, table.GetPlayerUnsafe(0))
}

