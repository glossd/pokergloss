package e2e

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/web/client/mq"
	"testing"
	"time"
)

func TestSeatReservationTimeout(t *testing.T) {
	timeOutSetUp(t)
	conf.Props.SeatReservationTimeout = time.Nanosecond

	table := InsertTable(t)

	restReserveSeat(t, table.ID.Hex(), 0)
	assertSeatReserved(t, 0)

	assertMessage(t, 1, func(as []*Asserter) {
		as[0].assertSeatReservedTimeout(0)
	})
}

func timeOutSetUp(t *testing.T) {
	prevPropsSetup(t)
	mq.IsTimeoutTestMQEnabled = true
	t.Cleanup(func() {
		mq.IsTimeoutTestMQEnabled = false
	})
}

func TestBuyIn_After_SeatReservationTimeout(t *testing.T) {
	timeOutSetUp(t)
	conf.Props.SeatReservationTimeout = time.Nanosecond

	table := InsertTable(t)

	restReserveSeat(t, table.ID.Hex(), 0)
	assertSeatReserved(t, 0)

	assertMessage(t, 1, func(as []*Asserter) {
		as[0].assertSeatReservedTimeout(0)
	})

	restBuyInStatus(t, table.ID.Hex(), 0, 400)
}

func TestBuyIn_Before_SeatReservationTimeout(t *testing.T) {
	timeOutSetUp(t)
	conf.Props.SeatReservationTimeout = 50 * time.Millisecond
	table := InsertTable(t)

	restReserveSeat(t, table.ID.Hex(), 0)
	assertSeatReserved(t, 0)
	restBuyIn(t, table.ID.Hex(), 0)
	assertBankroll(t, 0)

	// checking that the next event will be not a seatReservationTimeout
	time.Sleep(conf.Props.SeatReservationTimeout)
	restReserveSeat(t, table.ID.Hex(), 1, getToken(1))
	assertSeatReserved(t, 1)
}
