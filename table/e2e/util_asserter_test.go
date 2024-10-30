package e2e

import (
	"fmt"
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/services/model"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"testing"
)

type Asserter struct {
	T     *testing.T
	Event *mqws.Event
}

func NewAsserter(t *testing.T, event *mqws.Event) *Asserter {
	return &Asserter{T: t, Event: event}
}

func (a *Asserter) assertPayload(path string, expected interface{}, msg ...string) {
	assert.EqualValues(a.T, expected, gjson.Get(a.Event.Payload, path).Value(), msg)
}

func (a *Asserter) assertType(tet events.TET) {
	assert.EqualValues(a.T, tet, a.Event.Type)
}

func (a *Asserter) assertNotNilPayload(path string) {
	assert.NotNil(a.T, gjson.Get(a.Event.Payload, path).Value())
}

func (a *Asserter) assertNotEmptyPayload(path string) {
	assert.NotEmpty(a.T, gjson.Get(a.Event.Payload, path).Value())
}

func (a *Asserter) assertNotZeroPayload(path string) {
	v := gjson.Get(a.Event.Payload, path).Value()
	assert.NotNil(a.T, v)
	assert.NotZero(a.T, v)
}

func (a *Asserter) assertUndefinedPayload(path string) {
	assert.False(a.T, gjson.Get(a.Event.Payload, path).Exists())
}

func (a *Asserter) assertSeatReserved(position int) {
	a.assertType("seatReserved")
	a.assertPayload("table.seats.#", 1)
	a.assertPayload("table.seats.0.position", position)
}

func (a *Asserter) assertSeatReservedTimeout(position int) {
	a.assertType("seatReservationTimeout")
	a.assertPayload("table.seats.#", 1)
	a.assertPayload("table.seats.0.position", position)
	a.assertPayload("table.seats.0.player", nil)
}

// https://github.com/tidwall/gjson/pull/44#issuecomment-336664265
func (a *Asserter) assertNullPayload(path string) {
	result := gjson.Get(a.Event.Payload, path)
	assert.True(a.T, result.Type == gjson.Null && result.Exists(), fmt.Sprintf("%s is not null", path))
}

func (a *Asserter) assertBankroll(position int) {
	a.assertType("bankroll")
	a.assertPayload("table.seats.#", 1)
	a.assertPayload("table.seats.0.position", position)
	a.assertNotNilPayload("table.seats.0.player")
}

func (a *Asserter) assertAddChips(position int) {
	a.assertType("addChips")
	a.assertPayload("table.seats.#", 1)
	a.assertPayload("table.seats.0.position", position)
	a.assertNotNilPayload("table.seats.0.player")
}

func (a *Asserter) assertBettingAction(position int, action domain.ActionType) {
	a.assertType("playerAction")
	a.assertPayload("table.seats.#", 1)
	a.assertPayload("table.seats.0.position", position)
	a.assertPayload("table.seats.0.player.lastGameAction", action)
}

func (a *Asserter) assertTimeToDecide(position int) {
	a.assertType("timeToDecide")
	a.assertPayload("table.seats.#", 1)
	a.assertPayload("table.seats.0.position", position)
	a.assertNotZeroPayload("table.seats.0.player.timeoutAt")
}

func (a *Asserter) assertTimeToDecideTimeout(position int) {
	a.assertType("timeToDecideTimeout")
	a.assertPayload("table.seats.#", 1)
	a.assertPayload("table.seats.0.position", position)
	a.assertPayload("table.seats.0.player.lastGameAction", domain.Fold)
}

func (a *Asserter) assertReset(seatsNum int) {
	a.assertType("reset")
	a.assertPayload("table.pots.#", 0)
	a.assertPayload("table.communityCards.#", 0)
	if seatsNum > 0 {
		a.assertPayload("table.seats.#", seatsNum)
	} else {
		// todo maybe should return 0 seats
		a.assertUndefinedPayload("table.seats")
	}
}

func (a *Asserter) assertBlinds(bbPosition int, sbPosition int) {
	a.assertType("blinds")
	a.assertPayload("table.seats.#", 2)
	a.assertPayload("table.seats.0.position", bbPosition)
	a.assertPayload("table.seats.0.player.blind", "bigBlind")
	a.assertPayload("table.seats.1.position", sbPosition)
	a.assertPayload("table.seats.1.player.blind", "dealerSmallBlind")
}

func (a *Asserter) assertHoleCards(bbPosition int, sbPosition int) {
	a.assertType("holeCards")
	a.assertPayload("table.seats.#", 2)
	a.assertPayload("table.seats.0.player.cards.#", 2)
	a.assertPayload("table.seats.1.player.cards.#", 2)
}

func (a *Asserter) assertStackOverFlowPlayer() {
	a.assertType(events.StackOverflowPlayer)
}

func (a *Asserter) assertShowdown(position int, mucked bool) {
	a.assertType("showDown")
	a.assertPayload("table.seats.#", 1)
	a.assertPayload("table.seats.0.position", position)
	if mucked {
		a.assertUndefinedPayload("table.seats.0.player.cards")
	} else {
		a.assertPayload("table.seats.0.player.cards.#", 2)
	}
}

func (a *Asserter) assertWinners(winnersCount int) {
	a.assertType(events.Winners)
	a.assertPayload("table.winners.#", winnersCount)
}

func (a *Asserter) assertWinnersSimple() {
	a.assertType("winners")
	a.assertNotZeroPayload("table.winners.#")
}

func (a *Asserter) assertWinnersV2(winners []domain.Winner) {
	a.assertType(events.Winners)
	a.assertPayload("table.winners.#", len(winners))
	for i, winner := range winners {
		a.assertPayload(fmt.Sprintf("table.winners.%d.position", i), winner.Position)
		a.assertPayload(fmt.Sprintf("table.winners.%d.chips", i), winner.Chips)
	}
}

func (a *Asserter) assertIntent(intent *model.Intent) {
	a.assertType("intent")
	a.assertPayload("table.seats.#", 1)
	if intent == nil {
		a.assertNullPayload("table.seats.0.player.intent")
	} else {
		a.assertPayload("table.seats.0.player.intent.type", intent.Type)
		a.assertPayload("table.seats.0.player.intent.chips", intent.Chips)
	}
}

func (a *Asserter) assertNewBettingRound(roundType domain.RoundType) {
	a.assertType("newBettingRound")
	a.assertNotZeroPayload("table.pots.#")
	a.assertNotZeroPayload("table.pots.0.chips")
	a.assertPayload("roundType", roundType)
	if roundType == domain.FlopRound {
		a.assertPayload("newCards.#", 3)
	} else {
		a.assertPayload("newCards.#", 1)
	}
}

func (a *Asserter) assertPlayerLeft(position int) {
	a.assertType("playerLeft")
	a.assertPayload("table.seats.#", 1)
	a.assertPayload("table.seats.0.position", position)
}

func (a *Asserter) assertSitBack(position int) {
	a.assertType("sitBack")
	a.assertPayload("table.seats.#", 1)
	a.assertPayload("table.seats.0.position", position)
}

func (a *Asserter) assertMultiPlayersUpdate(tableId string, playersCount int) {
	a.assertType("multiPlayersUpdate")
	a.assertPayload("tableId", tableId)
	a.assertPayload("players.#", playersCount)
}

func (a *Asserter) assertMultiPlusPlayersUpdate(tableId string, playersCount int) {
	a.assertType("multiPlusPlayersUpdate")
	a.assertPayload("tableId", tableId)
	a.assertPayload("players.#", playersCount)
}

func (a *Asserter) assertMultiPlayerMove(toTableId string) {
	a.assertType("multiPlayerMove")
	a.assertPayload("tableId", toTableId)
}

func (a *Asserter) assertMultiLobby(lobby *domain.LobbyMulti) {
	a.assertType("multiLobby")
	a.assertNotNilPayload("lobby")
	a.assertPayload("lobby.tables.#", len(lobby.TableIDs))
}

func (a *Asserter) assertPlayerLeftPrize(position int, place int, prize int64) {
	a.assertType("playerLeft")
	a.assertPayload("leftPlayer.position", position)
	a.assertPayload("leftPlayer.tournamentInfo.place", place)
	a.assertPayload("leftPlayer.tournamentInfo.prize", prize)
}

func (a *Asserter) assertSitngoRegister(position int) {
	a.assertType("sitngoRegister")
	a.assertPayload("position", position)
}

func (a *Asserter) assertSitngoStart() {
	a.assertType("sitngoGameStart")
	a.assertNotEmptyPayload("table.id")
}

func (a *Asserter) assertMultiRegister() {
	a.assertType("multiRegister")
}

func (a *Asserter) assertMultiGameStart(tableId string) {
	a.assertType("multiGameStart")
	a.assertPayload("tableId", tableId)
}
