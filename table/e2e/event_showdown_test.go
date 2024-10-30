package e2e

import (
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/web/client/mq"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"testing"
)

func TestShowDown(t *testing.T) {
	t.Cleanup(cleanUp)
	domain.Algo = Algo_2P_MockZeroPositionLoses(t)

	tableID := RestCreatedTableRiverNext(t, 0, 1)

	restSetAutoMuckConfig(t, tableID.Hex(), 0, false)
	restMakeAction(t, tableID.Hex(), domain.Check, 0)
	readMessage()

	// River
	restMakeAction(t, tableID.Hex(), domain.Check, 1, secondPlayerToken)
	readMessage()

	restMakeAction(t, tableID.Hex(), domain.Check, 0)

	assertMessage(t, 3, func(as []*Asserter) {
		as[0].assertBettingAction(0, domain.Check)
		as[1].assertShowdown(1, false)
		as[2].assertTimeToDecide(0)
		assert.EqualValues(t, domain.ShowdownTable, gjson.Get(as[2].Event.Payload, "table.status").String())
	})
	restMakeShowDownAction(t, tableID.Hex(), domain.Muck, 0)

	assertMessage(t, 2, func(as []*Asserter) {
		as[0].assertShowdown(0, true)
		as[1].assertWinnersSimple()
	})
}

// https://youtu.be/w_WsO39YRx4?t=287
func TestStackOverflowPlayer(t *testing.T) {
	multiSetUp(t)
	domain.Algo = &domain.MockAlgo{}

	lobby := insertFullLobbyMulti(t, NewLobbyMultiParams{tableSize: 6, numOfUsers: 4})

	tableParams := domain.NewTableParams{
		Name:            "1",
		Size:            lobby.TableSize,
		BigBlind:        400,
		DecisionTimeout: lobby.DecisionTimeout,
		Identity:        authid.Identity{UserId: "Automatic"},
	}
	domain.Algo = domain.NewMockFull(0,
		domain.CardsStr(
			"2s", "7d",
			"Qh", "Ts",
			"8d", "Kh",
			"Js", "Qs",
		),
		domain.CardsStr("5c", "8s", "Ks", "Td", "5s"), // community cards
	)
	seats := lobby.IdensToSeats(lobby.Users)
	seats[0].Player.Stack = 1235
	seats[1].Player.Stack = 1769 + 200
	seats[2].Player.Stack = 387 + 400
	seats[3].Player.Stack = 574
	table, err := domain.NewTableMulti(lobby, tableParams, seats)
	assert.Nil(t, err)
	insertTable(t, table)
	hex := table.ID.Hex()

	restMakeAction(t, hex, domain.Call, 3, getToken(3))
	restMakeAction(t, hex, domain.Fold, 0, getToken(0))
	restMakeAction(t, hex, domain.Call, 1, getToken(1))
	restMakeAction(t, hex, domain.Check, 2, getToken(2))

	restMakeAction(t, hex, domain.AllIn, 1, getToken(1))
	restMakeAction(t, hex, domain.AllIn, 2, getToken(2))
	mq.ResetTestMQ()
	restMakeAction(t, hex, domain.AllIn, 3, getToken(3))
	assertMessage(t, 6, func(as []*Asserter) {
		as[0].assertType(events.PlayerMadeAction)

		as[1].assertType(events.StackOverflowPlayer)
		as[1].assertPayload("table.seats.0.player.stack", 1182)

		as[2].assertType(events.ShowDown)
		as[2].assertUndefinedPayload("table.seats.0.player.stack")
	})

}

func Algo_2P_MockZeroPositionLoses(t *testing.T) *domain.MockAlgo {
	mock, err := domain.NewMockAlgo(domain.CardsStr("2d", "7h", "As", "Ks", "Qs", "Js", "Ts", "8h", "4c"))
	assert.Nil(t, err)
	return mock
}

func Algo_2P_MockFirstPositionLoses(t *testing.T) *domain.MockAlgo {
	mock, err := domain.NewMockAlgo(domain.CardsStr("As", "Ks", "2d", "7h", "Qs", "Js", "Ts", "8h", "4c"))
	assert.Nil(t, err)
	return mock
}
