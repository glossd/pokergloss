package e2e

import (
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

// Had to optimize for low on resources gitlab runners

func TestDecisionTimeout(t *testing.T) {
	timeOutSetUp(t)

	sbPosition := 0
	bbPosition := 1
	table := NewTableWithStartedGame(t, sbPosition, bbPosition, -1)
	insertTable(t, table)

	restMakeAction(t, table.ID.Hex(), domain.Call, sbPosition, defaultToken)
	assertSimpleAction(t, sbPosition, bbPosition, domain.Call)
	res := actionhandler.DoDecisionTimeoutNoCtx(timeout.Key{TableID: table.ID, Position: bbPosition, Version: 1})
	assert.False(t, res)
	assertTimeToDecideTimeoutAndWinners(t, bbPosition, sbPosition)
}

func TestDecisionTimeoutOnFirstPlayer(t *testing.T) {
	timeOutSetUp(t)

	sbPosition := 0
	bbPosition := 1
	RestCreatedTableWithStartedGameTimeout(t, sbPosition, bbPosition, 5*time.Millisecond)

	assertTimeToDecideTimeoutAndWinners(t, sbPosition, bbPosition)
}

func TestDecisionTimeout_Before_Action(t *testing.T) {
	timeOutSetUp(t)

	sbPosition := 0
	bbPosition := 1
	table := NewTableWithStartedGame(t, sbPosition, bbPosition, time.Nanosecond)
	insertTable(t, table)

	restMakeAction(t, table.ID.Hex(), domain.Call, sbPosition, defaultToken)
	assertSimpleAction(t, sbPosition, bbPosition, domain.Call)

	assertTimeToDecideTimeoutAndWinners(t, bbPosition, sbPosition)

	restMakeActionStatus(t, table.ID.Hex(), domain.Check, bbPosition, http.StatusBadRequest, getToken(bbPosition))
}

func TestDecisionTimeout_After_Action(t *testing.T) {
	timeOutSetUp(t)

	sbPosition := 0
	bbPosition := 1
	table := NewTableWithStartedGame(t, sbPosition, bbPosition, 100*time.Millisecond)
	insertTable(t, table)

	restMakeAction(t, table.ID.Hex(), domain.Call, sbPosition, defaultToken)
	assertSimpleAction(t, sbPosition, bbPosition, domain.Call)

	restMakeAction(t, table.ID.Hex(), domain.Check, bbPosition, getToken(bbPosition))

	time.Sleep(table.DecisionTimeout + 5*time.Millisecond)
	assertActionNewBettingRound(t, bbPosition, domain.Check, domain.FlopRound, bbPosition)
}

func TestDecisionTimeout_On_AnotherBuyIn(t *testing.T) {
	timeOutSetUp(t)

	sbPosition := 0
	bbPosition := 1
	table := NewTableWithStartedGame(t, sbPosition, bbPosition, 100*time.Millisecond)
	insertTable(t, table)

	restMakeAction(t, table.ID.Hex(), domain.Call, sbPosition, defaultToken)
	assertSimpleAction(t, sbPosition, bbPosition, domain.Call)

	restReserveSeat(t, table.ID.Hex(), 2, getToken(2))
	assertSeatReserved(t, 2)
	restBuyIn(t, table.ID.Hex(), 2, getToken(2))
	assertBankroll(t, 2)

	assertTimeToDecideTimeoutAndWinners(t, bbPosition, sbPosition)
}

func TestDecisionTimeoutEmptyCommunityCards(t *testing.T) {
	timeOutSetUp(t)

	sbPosition := 0
	bbPosition := 1
	table := NewTableWithStartedGame(t, sbPosition, bbPosition, -1)
	err := table.MakeAction(sbPosition, getIden(sbPosition), domain.CallAction)
	assert.Nil(t, err)
	insertTable(t, table)
	tableID := table.ID

	restMakeAction(t, tableID.Hex(), domain.Check, bbPosition, secondPlayerToken)

	assertActionNewBettingRound(t, bbPosition, domain.Check, domain.FlopRound, bbPosition)

	actionhandler.DoDecisionTimeoutNoCtx(timeout.Key{TableID: tableID, Position: bbPosition, Version: 1})
	assertTimeToDecideTimeoutAndWinners(t, bbPosition, sbPosition)

	assertReset(t, bbPosition, sbPosition)

	table, err = db.FindTableNoCtx(tableID)
	assert.Nil(t, err)
	assert.EqualValues(t, domain.CommunityCards{}, *table.CommunityCards)
}

func NewTableWithStartedGame(t *testing.T, defaultPlayerPosition, secondPlayerPosition int, timeout time.Duration) *domain.Table {
	table := NewTableTimeout(t, timeout)

	err := table.ReserveSeat(defaultPlayerPosition, defaultIdentity)
	assert.Nil(t, err)
	_, err = table.BuyIn(250, defaultPlayerPosition, defaultIdentity)
	assert.Nil(t, err)

	err = table.ReserveSeat(secondPlayerPosition, secondIdentity)
	assert.Nil(t, err)
	_, err = table.BuyIn(250, secondPlayerPosition, secondIdentity)
	assert.Nil(t, err)

	return table
}
