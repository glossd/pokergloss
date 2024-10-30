package e2e

import (
	"testing"
)

func TestBankroll(t *testing.T) {
	t.Cleanup(cleanUp)
	table := InsertTable(t)

	position := 0
	restReserveSeat(t, table.ID.Hex(), position)
	assertSeatReserved(t, position)

	restBuyIn(t, table.ID.Hex(), position)
	assertBankroll(t, position)
}

func assertBankroll(t *testing.T, position int) {
	assertMessage(t, 1, func(as []*Asserter) {
		as[0].assertBankroll(position)
	})
}
