package e2e

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	"github.com/glossd/pokergloss/table/web/client/mq"
	"github.com/glossd/pokergloss/table/web/client/mqpub"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestFull2PlayerGame(t *testing.T) {
	//#########################################################
	//# This test shouldn't have any errors or warning in log #
	//#########################################################
	t.Cleanup(cleanUp)
	domain.Algo = Algo_2P_MockZeroPositionLoses(t)

	table := InsertTable(t)
	tableID := table.ID

	sbPosition := 0
	bbPosition := 1

	restReserveSeat(t, tableID.Hex(), sbPosition, defaultToken)
	assertSeatReserved(t, sbPosition)
	restBuyIn(t, tableID.Hex(), sbPosition, defaultToken)
	assertBankroll(t, sbPosition)

	restReserveSeat(t, tableID.Hex(), bbPosition, secondPlayerToken)
	assertSeatReserved(t, bbPosition)
	restBuyIn(t, tableID.Hex(), bbPosition, secondPlayerToken)
	assertBankroll(t, bbPosition)

	assertStartHand(t, bbPosition, sbPosition, sbPosition)

	restMakeAction(t, tableID.Hex(), domain.Call, sbPosition, defaultToken)
	assertSimpleAction(t, sbPosition, bbPosition, domain.Call)

	restMakeAction(t, tableID.Hex(), domain.Check, bbPosition, secondPlayerToken)
	assertActionNewBettingRound(t, bbPosition, domain.Check, domain.FlopRound, bbPosition)

	restMakeAction(t, tableID.Hex(), domain.Check, bbPosition, secondPlayerToken)
	assertSimpleAction(t, bbPosition, sbPosition, domain.Check)

	restMakeAction(t, tableID.Hex(), domain.Check, sbPosition, defaultToken)
	assertActionNewBettingRound(t, sbPosition, domain.Check, domain.TurnRound, bbPosition)

	restMakeAction(t, tableID.Hex(), domain.Check, bbPosition, secondPlayerToken)
	assertSimpleAction(t, bbPosition, sbPosition, domain.Check)

	restMakeAction(t, tableID.Hex(), domain.Check, sbPosition, defaultToken)
	assertActionNewBettingRound(t, sbPosition, domain.Check, domain.RiverRound, bbPosition)

	// River actions

	restMakeAction(t, tableID.Hex(), domain.Check, bbPosition, secondPlayerToken)
	assertSimpleAction(t, bbPosition, sbPosition, domain.Check)

	restMakeAction(t, tableID.Hex(), domain.Check, sbPosition, defaultToken)
	assertMessage(t, 4, func(as []*Asserter) {
		as[0].assertBettingAction(sbPosition, domain.Check)
		as[1].assertShowdown(bbPosition, false)
		as[2].assertShowdown(sbPosition, true)
		as[3].assertWinners(1)
	})

	gameEnd := <-mq.TestGameEndMQ
	assert.EqualValues(t, 1, len(gameEnd.Winners))
	assert.EqualValues(t, 4, gameEnd.Winners[0].Chips)
	assert.EqualValues(t, "Straight Flush", gameEnd.Winners[0].Hand)
	assert.EqualValues(t, 2, len(gameEnd.Players))
	assert.EqualValues(t, 2, gameEnd.Players[0].WageredChips)
	assert.EqualValues(t, 2, gameEnd.Players[1].WageredChips)
	assert.EqualValues(t, 5, len(gameEnd.CommunityCards))

	newSbPosition := bbPosition
	newBbPosition := sbPosition

	assertStartHand(t, newBbPosition, newSbPosition, newSbPosition)
}

func TestFull3PlayerGame(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.Table.GameEndMinTimeout = 0

	table := InsertTable(t)
	tableID := table.ID
	tableIdHex := tableID.Hex()

	fPosition := 0
	sPosition := 1
	tPosition := 2

	restReserveSeat(t, tableIdHex, fPosition)
	restBuyIn(t, tableIdHex, fPosition)

	restReserveSeat(t, tableIdHex, sPosition, secondPlayerToken)
	restBuyIn(t, tableIdHex, sPosition, secondPlayerToken)

	restReserveSeat(t, tableIdHex, tPosition, thirdPlayerToken)
	restBuyIn(t, tableIdHex, tPosition, thirdPlayerToken)

	restMakeAction(t, tableIdHex, domain.Fold, fPosition)

	// New game, give it time
	time.Sleep(10 * time.Millisecond)

	bbPosition := 0
	dealerPosition := 1
	sbPosition := 2

	restMakeAction(t, tableIdHex, domain.Fold, dealerPosition, secondPlayerToken)
	restMakeAction(t, tableIdHex, domain.Call, sbPosition, thirdPlayerToken)
	restMakeAction(t, tableIdHex, domain.Check, bbPosition)
	restMakeAction(t, tableIdHex, domain.Check, sbPosition, thirdPlayerToken)
	restMakeAction(t, tableIdHex, domain.Check, bbPosition)
}

func TestBigBlindForcedAllIn(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.GameEndMinTimeout = -1
	mq.IsTimeoutTestMQEnabled = true

	tableID := RestCreatedTableWithStartedGame(t, 0, 1)
	hex := tableID.Hex()
	restMakeBetAction(t, hex, domain.Raise, 248, 0, getToken(0))
	restMakeAction(t, hex, domain.Call, 1, getToken(1))
	restMakeAction(t, hex, domain.Check, 1, getToken(1))

	mq.ResetTestMQ()
	restMakeAction(t, hex, domain.Fold, 0, getToken(0))

	assertMessage(t, 3, func(as []*Asserter) {
		as[0].assertBettingAction(0, domain.Fold)
		as[1].assertShowdown(1, true)
		as[2].assertWinners(1)
	})

	domain.Algo = Algo_2P_MockFirstPositionLoses(t)
	startGameManually(t, tableID)
	msg := readMessage()

	assert.Len(t, msg.UserEvents.BeforeEvents, 2)
	NewAsserter(t, msg.UserEvents.BeforeEvents[0]).assertReset(9)
	assert.EqualValues(t, 1, gjson.Get(msg.UserEvents.BeforeEvents[0].Payload, "table.seats.0.player.stack").Int())
	NewAsserter(t, msg.UserEvents.BeforeEvents[1]).assertBlinds(0, 1)
	assert.EqualValues(t, 0, gjson.Get(msg.UserEvents.BeforeEvents[1].Payload, "table.seats.0.player.stack").Int())
	assert.EqualValues(t, 1, gjson.Get(msg.UserEvents.BeforeEvents[1].Payload, "table.seats.0.player.totalRoundBet").Int())

	assert.Len(t, msg.UserEvents.NotFoundUsersEvents, 1)
	NewAsserter(t, msg.UserEvents.NotFoundUsersEvents[0]).assertHoleCards(0, 1)
	assert.False(t, gjson.Get(msg.UserEvents.NotFoundUsersEvents[0].Payload, "table.seats.0.player.stack").Exists())
	assert.EqualValues(t, 1, gjson.Get(msg.UserEvents.NotFoundUsersEvents[0].Payload, "table.seats.0.player.totalRoundBet").Int())

	assert.Len(t, msg.UserEvents.AfterEvents, 2)
	assert.EqualValues(t, events.NewBettingRound, msg.UserEvents.AfterEvents[0].Type)
	assert.EqualValues(t, 0, gjson.Get(msg.UserEvents.AfterEvents[0].Payload, "table.seats.0.player.stack").Int())
	assert.EqualValues(t, 0, gjson.Get(msg.UserEvents.AfterEvents[0].Payload, "table.seats.0.player.totalRoundBet").Int())

	assert.EqualValues(t, events.Winners, msg.UserEvents.AfterEvents[1].Type)
	assert.EqualValues(t, 0, gjson.Get(msg.UserEvents.AfterEvents[1].Payload, "table.seats.0.player.stack").Int())
	assert.EqualValues(t, 0, gjson.Get(msg.UserEvents.AfterEvents[1].Payload, "table.seats.0.player.totalRoundBet").Int())
}

func TestEveryoneAllIn_ShouldReturnAccurateReset(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.GameEndMinTimeout = -1
	mq.IsTimeoutTestMQEnabled = true
	domain.Algo = Algo_2P_MockZeroPositionLoses(t)
	tableID := RestCreatedTableWithStartedGame(t, 0, 1)
	hex := tableID.Hex()
	restMakeAction(t, hex, domain.AllIn, 0, getToken(0))
	restMakeAction(t, hex, domain.AllIn, 1, getToken(1))
	mq.ResetTestMQ()
	startGameManually(t, tableID)
	assertMessage(t, 1, func(as []*Asserter) {
		a := as[0]
		a.assertType(events.Reset)
		a.assertPayload("table.seats.0.player.stack", 0)
	})
}

func startGameManually(t *testing.T, tableID primitive.ObjectID) {
	mqpub.PublishTimeoutEvent(&timeout.Event{Type: timeout.StartGame, At: timeutil.Now(), Key: timeout.Key{TableID: tableID, Version: findTable(t, tableID).GameFlowVersion}})
}
