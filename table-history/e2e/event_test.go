package e2e

import (
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/table-history/db"
	"github.com/glossd/pokergloss/table-history/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommon(t *testing.T) {
	t.Cleanup(cleanUpDB)
	events := domain.ToEvents(&mqws.TableMessage{
		ToEntityIds: []string{"123"},
		UserEvents:  nil,
		Events: []*mqws.Event{
			{Type: "playerAction", Payload: `{"table": {"seats":[{"position":0},{"position":1}]}}`},
		},
	})
	assert.EqualValues(t, 1, len(events))
	assert.Nil(t, db.InsertManyEvents(events))
}

// I sorted by created at, should've used _id for sorting
func TestEventsOrder(t *testing.T) {
	t.Cleanup(cleanUpDB)
	events := domain.ToEvents(&mqws.TableMessage{
		ToEntityIds: []string{"123"},
		UserEvents:  nil,
		Events: []*mqws.Event{
			{Type: "playerAction", Payload: `{"table":{"decidingPosition":-1,"maxRoundBet":0,"seats":[{"blind":"bigBlind","player":{"blind":"bigBlind","intent":null,"isDeciding":true,"lastGameAction":"check","lastGameBet":2,"picture":"https://storage.googleapis.com/pokerblow-avatars/cLhhrCC0sMUdpIHGlngodn1lKLl2","position":2,"stack":98,"status":"playing","totalRoundBet":2,"userId":"cLhhrCC0sMUdpIHGlngodn1lKLl2","username":"Sea_Man"},"position":2}],"totalPot":4}}}`},
			{Type: "newBettingRound", Payload: `{"newCards":["Qd","5c","5h"],"roundType":"flop","table":{"communityCards":["Qd","5c","5h"],"pot":4,"pots":[{"chips":4,"winnerPositions":null}],"seats":[{"blind":"dealerSmallBlind","player":{"intent":null,"isDeciding":false,"position":1,"totalRoundBet":0},"position":1},{"blind":"bigBlind","player":{"intent":null,"isDeciding":false,"position":2,"totalRoundBet":0},"position":2}],"totalPot":4}}}`},
			{Type: "timeToDecide", Payload: `{"table":{"decidingPosition":2,"lastAggressorPosition":-1,"maxRoundBet":0,"seats":[{"blind":"bigBlind","player":{"intent":null,"isDeciding":true,"position":2,"timeoutAt":1610368453443},"position":2}]}}}`},
		},
	})

	assert.Nil(t, db.InsertManyEvents(events))

	events2 := domain.ToEvents(&mqws.TableMessage{
		ToEntityIds: []string{"123"},
		UserEvents:  nil,
		Events: []*mqws.Event{
			{Type: "playerAction", Payload: `{"table":{"decidingPosition":-1,"maxRoundBet":0,"seats":[{"blind":"bigBlind","player":{"blind":"bigBlind","intent":null,"isDeciding":false,"lastGameAction":"check","lastGameBet":0,"picture":"https://storage.googleapis.com/pokerblow-avatars/bmO68rNtnjfziNkeQH41kPboXi83","position":1,"stack":250,"status":"playing","totalRoundBet":0,"userId":"bmO68rNtnjfziNkeQH41kPboXi83","username":"eminem"},"position":1}],"totalPot":4}}}`},
			{Type: "timeToDecide", Payload: `{"table":{"decidingPosition":2,"lastAggressorPosition":-1,"maxRoundBet":0,"seats":[{"blind":"dealerSmallBlind","player":{"intent":null,"isDeciding":true,"position":2,"timeoutAt":1610368755316},"position":2}]}}`},
		},
	})

	assert.Nil(t, db.InsertManyEvents(events2))

	allEvents, err := db.FindAll()
	assert.Nil(t, err)

	assert.EqualValues(t, 5, len(allEvents))
	assert.EqualValues(t, "playerAction", allEvents[0].Type)
	assert.EqualValues(t, "newBettingRound", allEvents[1].Type)
	assert.EqualValues(t, "timeToDecide", allEvents[2].Type)
	assert.EqualValues(t, "playerAction", allEvents[3].Type)
	assert.EqualValues(t, "timeToDecide", allEvents[4].Type)
}
