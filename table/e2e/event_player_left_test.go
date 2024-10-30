package e2e

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayerLeft(t *testing.T) {
	t.Cleanup(cleanUp)

	tablePosition := 0
	table := NewTable(t)
	err := table.ReserveSeat(tablePosition, defaultIdentity)
	assert.Nil(t, err)
	_, err = table.BuyIn(250, tablePosition, defaultIdentity)
	assert.Nil(t, err)

	insertTable(t, table)

	restPlayerStand(t, table.ID, tablePosition)

	assertPlayerLeft(t, tablePosition)

	tableSaved, err := db.FindTableNoCtx(table.ID)
	assert.Nil(t, err)
	assert.Nil(t, tableSaved.Seats[tablePosition].Player)
}

func TestPlayerLeftWhilePlaying(t *testing.T) {
	t.Cleanup(cleanUp)

	sbPosition := 0
	bbPosition := 1
	tableID := RestCreatedTableWithStartedGame(t, sbPosition, bbPosition)

	restPlayerStand(t, tableID, bbPosition, secondPlayerToken)

	restMakeAction(t, tableID.Hex(), domain.Call, sbPosition)
	assertAction_Leave_Winners(t, sbPosition, domain.Call, bbPosition)
}

func TestPlayerLeft_AndGameEnd(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.Table.GameEndMinTimeout = 0

	tableID := RestCreatedTableTurnNext(t, 0, 1)

	restPlayerStand(t, tableID, 1, secondPlayerToken)

	restMakeAction(t, tableID.Hex(), domain.Fold, 0)

	assertMessage(t, 3, func(as []*Asserter) {
		as[0].assertBettingAction(0, domain.Fold)
		as[1].assertShowdown(1, true)
		as[2].assertWinnersV2([]domain.Winner{{Position: 1, Chips: 4}})
	})

	savedTable, err := db.FindTableNoCtx(tableID)
	assert.Nil(t, err)
	if savedTable.IsWaiting() {
		assert.True(t, savedTable.GetSeatUnsafe(1).IsFree())
	}
}

func TestPlayerLeft_CancelReservation(t *testing.T) {
	t.Cleanup(cleanUp)

	table := InsertTable(t)

	restReserveSeat(t, table.ID.Hex(), 0)
	assertSeatReserved(t, 0)

	restCancelSeatReservation(t, table.ID.Hex(), 0)
	assertPlayerLeft(t, 0)
}

func TestPlayerLeft_WhileDeciding(t *testing.T) {
	t.Cleanup(cleanUp)

	tableID := RestCreatedTableTurnNext(t, 0, 1)

	restPlayerStand(t, tableID, 0, getToken(0))

	assertMessage(t, 4, func(as []*Asserter) {
		as[0].assertBettingAction(0, domain.Fold)
		as[1].assertShowdown(1, true)
		as[2].assertWinners(1)
		as[3].assertPlayerLeft(0)
	})
}

func assertPlayerLeft(t *testing.T, position int) {
	assertMessage(t, 1, func(as []*Asserter) {
		as[0].assertPlayerLeft(position)
	})
}
