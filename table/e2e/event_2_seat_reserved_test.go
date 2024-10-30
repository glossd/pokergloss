package e2e

import (
	"testing"
)

func TestReserveSeat(t *testing.T) {
	t.Cleanup(cleanUp)
	table := InsertTable(t)

	restReserveSeat(t, table.ID.Hex(), 0)
	assertSeatReserved(t, 0)
}

func assertSeatReserved(t *testing.T, position int) {
	assertMessage(t, 1, func(as []*Asserter) {
		as[0].assertSeatReserved(position)
	})
}
