package e2e

import (
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

// Player with lowest position is dealer. (here is only two players to dealer is smallBlind too)
// according to the mock algo domain.Algo#ChooseDealer

const defaultBuyIn = 250
const defaultBigBlind = 2
const defaultSmallBlind = 1

func RestCreatedTableWithStartedGame(t *testing.T, firstPosition, secondPosition int) primitive.ObjectID {
	return RestCreatedTableWithStartedGameTimeout(t, firstPosition, secondPosition, -1)
}

func RestCreatedTableWithStartedGameTimeout(t *testing.T, firstPosition, secondPosition int, timeout time.Duration) primitive.ObjectID {
	table := InsertTableTimeout(t, timeout)

	reserveAndBuyIn(t, firstPosition, table, defaultToken)
	reserveAndBuyIn(t, secondPosition, table, secondPlayerToken)

	table, err := db.FindTableNoCtx(table.ID)
	assert.Nil(t, err)
	if table.Status == domain.PlayingTable {
		assertStartHand(t, table.BigBlindPosition(), table.SmallBlindPosition(), table.DecidingPosition)
	}
	if table.Status == domain.WaitingTable {
		readMessage()
	}

	return table.ID
}

func reserveAndBuyIn(t *testing.T, pos int, table *domain.Table, token string) {
	restReserveSeat(t, table.ID.Hex(), pos, token)
	assertSeatReserved(t, pos)
	restBuyIn(t, table.ID.Hex(), pos, token)
	assertBankroll(t, pos)
}

func reserveAndBuyInNoAsserts(t *testing.T, pos int, table *domain.Table, token string) {
	restReserveSeat(t, table.ID.Hex(), pos, token)
	restBuyIn(t, table.ID.Hex(), pos, token)
}

func RestCreatedTableFlopNext(t *testing.T, sbPosition, bbPosition int) primitive.ObjectID {
	tableID := RestCreatedTableWithStartedGame(t, sbPosition, bbPosition)
	restMakeAction(t, tableID.Hex(), domain.Call, sbPosition)
	readMessage()
	return tableID
}

func RestCreatedTableTurnNext(t *testing.T, sbPosition, bbPosition int) primitive.ObjectID {
	tableID := RestCreatedTableFlopNext(t, sbPosition, bbPosition)
	restMakeAction(t, tableID.Hex(), domain.Check, bbPosition, secondPlayerToken)
	readMessage()
	// Flop
	restMakeAction(t, tableID.Hex(), domain.Check, bbPosition, secondPlayerToken)
	readMessage()
	return tableID
}

func RestCreatedTableRiverNext(t *testing.T, sbPosition, bbPosition int) primitive.ObjectID {
	tableID := RestCreatedTableTurnNext(t, sbPosition, bbPosition)
	restMakeAction(t, tableID.Hex(), domain.Check, sbPosition)
	readMessage()
	// Turn
	restMakeAction(t, tableID.Hex(), domain.Check, bbPosition, secondPlayerToken)
	readMessage()
	return tableID
}

func RestCreatedTableEndNext(t *testing.T, sbPosition, bbPosition int) primitive.ObjectID {
	tableID := RestCreatedTableRiverNext(t, sbPosition, bbPosition)
	restMakeAction(t, tableID.Hex(), domain.Check, sbPosition)
	readMessage()
	// River
	restMakeAction(t, tableID.Hex(), domain.Check, bbPosition, secondPlayerToken)
	readMessage()
	return tableID
}
