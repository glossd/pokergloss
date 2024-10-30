package e2e

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/player"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddChips(t *testing.T) {
	t.Cleanup(cleanUp)

	table := InsertTable(t)

	restReserveSeat(t, table.ID.Hex(), 0)
	restBuyIn(t, table.ID.Hex(), 0)
	restAddChips(t, table.ID.Hex(), 0)

	assertSeatReserved(t, 0)
	assertBankroll(t, 0)
	assertMessage(t, 1, func(as []*Asserter) {
		as[0].assertAddChips(0)
	})
}

func TestRaceAddChipsOnGameEnd(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.Table.GameEndMinTimeout = -1

	tableId := RestCreatedTableEndNext(t, 0, 1)
	restMakeAction(t, tableId.Hex(), domain.Check, 0, getToken(0))

	table := findTable(t, tableId)
	actionhandler.DoStartGameNoCtx(timeout.Key{
		TableID: table.ID,
		Version: table.GameFlowVersion,
	})

	assert.EqualValues(t, db.ErrVersionNotMatch, player.AddChipsOnTable(chipsParams(t, table.ID, 0), table))
}
