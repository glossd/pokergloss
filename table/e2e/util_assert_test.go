package e2e

import (
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/web/client/mq"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func assertMessage(t *testing.T, length int, asserters func(as []*Asserter)) {
	msg := readMessage()
	assert.NotZero(t, len(msg.ToEntityIds))
	assert.Nil(t, msg.UserEvents)
	assert.Len(t, msg.Events, length)
	var acc []*Asserter
	for _, event := range msg.Events {
		acc = append(acc, NewAsserter(t, event))
	}
	asserters(acc)
}

func readMessage() *mqws.TableMessage {
	msg := <-mq.TestMQ
	return msg
}

func readNewsMessage() *mqws.Message {
	msg := <-mq.TestNewsMQ
	return msg
}

func assertStartHand(t *testing.T, bbPos int, sbPos int, decidingPos int) {
	msg := readMessage()
	assert.Len(t, msg.UserEvents.BeforeEvents, 2)
	NewAsserter(t, msg.UserEvents.BeforeEvents[0]).assertReset(9)
	NewAsserter(t, msg.UserEvents.BeforeEvents[1]).assertBlinds(bbPos, sbPos)

	assert.Len(t, msg.UserEvents.NotFoundUsersEvents, 1)
	NewAsserter(t, msg.UserEvents.NotFoundUsersEvents[0]).assertHoleCards(bbPos, sbPos)

	// todo test real user events

	assert.Len(t, msg.UserEvents.AfterEvents, 1)
	NewAsserter(t, msg.UserEvents.AfterEvents[0]).assertTimeToDecide(decidingPos)
}

func assertStartHandTableSize(t *testing.T, bbPos, sbPos, decidingPos, tableSize int) {
	msg := readMessage()
	assert.Len(t, msg.UserEvents.BeforeEvents, 2)
	NewAsserter(t, msg.UserEvents.BeforeEvents[0]).assertReset(tableSize)
	NewAsserter(t, msg.UserEvents.BeforeEvents[1]).assertBlinds(bbPos, sbPos)

	assert.Len(t, msg.UserEvents.NotFoundUsersEvents, 1)
	NewAsserter(t, msg.UserEvents.NotFoundUsersEvents[0]).assertHoleCards(bbPos, sbPos)

	// todo test real user events

	assert.Len(t, msg.UserEvents.AfterEvents, 1)
	NewAsserter(t, msg.UserEvents.AfterEvents[0]).assertTimeToDecide(decidingPos)
}

func assertTimeToDecideTimeoutAndStackOverFlowAndWinners(t *testing.T, timeoutPosition int, winnerPosition int) {
	assertMessage(t, 4, func(as []*Asserter) {
		as[0].assertTimeToDecideTimeout(timeoutPosition)
		as[1].assertStackOverFlowPlayer()
		as[2].assertShowdown(winnerPosition, true)
		as[3].assertWinners(1)
	})
}

func assertTimeToDecideTimeoutAndWinners(t *testing.T, timeoutPosition int, winnerPosition int) {
	assertMessage(t, 3, func(as []*Asserter) {
		as[0].assertTimeToDecideTimeout(timeoutPosition)
		as[1].assertShowdown(winnerPosition, true)
		as[2].assertWinners(1)
	})
}

func assertActionNewBettingRound(t *testing.T, pos int, action domain.ActionType, roundType domain.RoundType, decidingPos int) {
	assertMessage(t, 3, func(as []*Asserter) {
		as[0].assertBettingAction(pos, action)
		as[1].assertNewBettingRound(roundType)
		as[2].assertTimeToDecide(decidingPos)
	})
}

func assertAction_Leave_Winners(t *testing.T, pos int, action domain.ActionType, leftPos int) {
	assertMessage(t, 5, func(as []*Asserter) {
		as[0].assertBettingAction(pos, action)
		as[1].assertBettingAction(leftPos, domain.Fold)
		as[2].assertShowdown(pos, true)
		as[3].assertWinners(1)
		as[4].assertPlayerLeft(leftPos)
	})
}

func assertReset(t *testing.T, bbPosition int, sbPosition int) {
	assertMessage(t, 1, func(as []*Asserter) {
		as[0].assertReset(9)
	})
}

func assertMultiGameStart(t *testing.T, tableID string) {
	msg := readNewsMessage()
	for _, event := range msg.Events {
		NewAsserter(t, event).assertMultiGameStart(tableID)
	}
}

func assertMultiGameStartSimple(t *testing.T) {
	msg := readNewsMessage()
	for _, event := range msg.Events {
		NewAsserter(t, event).assertType("multiGameStart")
	}
}

func assertTableDontExist(t *testing.T, id primitive.ObjectID) {
	_, err := db.FindTableNoCtx(id)
	assert.EqualValues(t, mongo.ErrNoDocuments, err)
}

func assertSitngoRegister(t *testing.T, position int) {
	assertMessage(t, 1, func(as []*Asserter) {
		as[0].assertSitngoRegister(position)
	})
}

func assertSitngoGameStart(t *testing.T, tableID string) {
	msg := readNewsMessage()
	assert.EqualValues(t, 1, len(msg.Events))
	assert.EqualValues(t, events.SitngoGameStartType, msg.Events[0].Type)
	assert.EqualValues(t, tableID, gjson.Get(msg.Events[0].Payload, "tableId").String())
}

func assertSitngoGameStartSimple(t *testing.T) {
	msg := readNewsMessage()
	assert.EqualValues(t, 1, len(msg.Events))
	assert.EqualValues(t, events.SitngoGameStartType, msg.Events[0].Type)
	assert.NotEmpty(t, gjson.Get(msg.Events[0].Payload, "tableId").String())
}
