package e2e

import (
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGameStart(t *testing.T) {
	t.Cleanup(cleanUp)

	table := InsertTable(t)

	sbPosition := 0 // default
	bbPosition := 2 // second

	restReserveSeat(t, table.ID.Hex(), sbPosition)
	assertSeatReserved(t, sbPosition)

	restBuyIn(t, table.ID.Hex(), sbPosition)
	assertBankroll(t, sbPosition)

	restReserveSeat(t, table.ID.Hex(), bbPosition, secondPlayerToken)
	assertSeatReserved(t, bbPosition)

	restBuyIn(t, table.ID.Hex(), bbPosition, secondPlayerToken)
	assertBankroll(t, bbPosition)

	assertStartHand(t, bbPosition, sbPosition, sbPosition)
}

func TestPlayerFieldsAfterGameStarted(t *testing.T) {
	t.Cleanup(cleanUp)
	domain.Algo = &domain.MockAlgo{}

	sbPosition := 0
	bbPosition := 1
	tableId := RestCreatedTableWithStartedGame(t, sbPosition, bbPosition)

	table, err := db.FindTableNoCtx(tableId)
	assert.Nil(t, err)
	assert.EqualValues(t, 249, table.GetPlayerUnsafe(sbPosition).Stack)
	assert.EqualValues(t, 1, table.GetPlayerUnsafe(sbPosition).TotalRoundBet)
	assert.EqualValues(t, 248, table.GetPlayerUnsafe(bbPosition).Stack)
	assert.EqualValues(t, 2, table.GetPlayerUnsafe(bbPosition).TotalRoundBet)
}
