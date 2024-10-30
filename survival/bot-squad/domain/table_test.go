package domain

import (
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"testing"
)

func TestNewTable(t *testing.T) {
	table := getEmptyTable(t)
	assert.EqualValues(t, "Really Long T", table.Name)
	assert.EqualValues(t, 6, len(table.Seats))
}

func getEmptyTable(t *testing.T) *Table {
	table, err := NewTable(emptyTable())
	assert.Nil(t, err)
	return table
}

func emptyTable() []byte {
	return []byte(`{"id":"60817bd8383a2d835945c6e9","name":"Really Long T","size":6,"bigBlind":2,"smallBlind":1,"decisionTimeoutSec":600,"bettingLimit":"NL","minBuyIn":100,"maxBuyIn":400,"type":"cashGame","occupied":0,"avgStake":0,"avgPot":0,"seats":[{"position":0,"blind":"","player":null},{"position":1,"blind":"","player":null},{"position":2,"blind":"","player":null},{"position":3,"blind":"","player":null},{"position":4,"blind":"","player":null},{"position":5,"blind":"","player":null}],"status":"waiting","maxRoundBet":0,"bettingLimitChips":9223372036854775807,"pots":null,"totalPot":0,"communityCards":[],"decidingPosition":-1,"lastAggressorPosition":-1}`)
}

// Bot has "7s", "6c"
// Board has "9c", "8c", "5s", "2d"
// Turn: bot checks, user bets 88, totalPot = 268, bot folds
func TestBug(t *testing.T) {
	initState := `{"id":"60b4d1186fc66047451c509e","name":"Mistral","size":2,"bigBlind":2,"smallBlind":1,"decisionTimeoutSec":60,"bettingLimit":"NL","minBuyIn":100,"maxBuyIn":400,"type":"sitAndGo","occupied":2,"avgStake":0,"avgPot":0,"seats":[{"position":0,"blind":"dealerSmallBlind","player":{"userId":"cLhhrCC0sMUdpIHGlngodn1lKLl2","username":"Sea_Man","picture":"https://storage.googleapis.com/pokerblow-avatars/cLhhrCC0sMUdpIHGlngodn1lKLl2","position":0,"stack":299,"status":"playing","cards":["Xx","Xx"],"blind":"dealerSmallBlind","totalRoundBet":1,"lastGameBet":1,"lastGameAction":"","timeoutAt":1622462805735,"intent":null,"level":6,"bankBalance":266365,"bankRank":3,"marketItemId":"pacifier","isDeciding":true}},{"position":1,"blind":"bigBlind","player":{"userId":"DIzWERaQmnZUxq9OBq4QXGigdpy2","username":"Wonder_Woman","picture":"https://storage.googleapis.com/pokerblow-avatars/DIzWERaQmnZUxq9OBq4QXGigdpy2","position":1,"stack":298,"status":"playing","cards":["Kh","8d"],"blind":"bigBlind","totalRoundBet":2,"lastGameBet":2,"lastGameAction":"","intent":null,"level":6,"bankBalance":276161,"bankRank":2,"marketItemId":"lightBulb"}}],"status":"playing","maxRoundBet":2,"bettingLimitChips":9223372036854775807,"pots":null,"totalPot":3,"communityCards":[],"decidingPosition":0,"lastAggressorPosition":-1,"tournamentAttrs":{"lobbyId":"60b4d10c039f12529aaeabfd","name":"Mistral","startAt":1622462744735,"levelIncreaseAt":1622462924735,"nextSmallBlind":2,"prizes":[{"place":1,"prize":600}]},"themeId":"","isSurvival":false,"survivalLevel":0}`
	table, err := NewTable([]byte(initState))
	assert.Nil(t, err)

	merge := func(data string) {
		events := extractEvents(data)
		for _, event := range events {
			assert.Nil(t, event.Merge(defaultConfig, table))
		}
	}

	resetHoleCards := `[{"type":"reset","payload":{"table":{"bettingLimitChips":0,"communityCards":[],"decidingPosition":-1,"lastAggressorPosition":-1,"maxRoundBet":0,"pots":[],"seats":[{"blind":"dealerSmallBlind","player":{"blind":"","cards":[],"intent":null,"lastGameAction":"","lastGameBet":0,"picture":"https://storage.googleapis.com/pokerblow-avatars/cLhhrCC0sMUdpIHGlngodn1lKLl2","position":0,"stack":304,"status":"playing","totalRoundBet":0,"userId":"cLhhrCC0sMUdpIHGlngodn1lKLl2","username":"Sea_Man"},"position":0},{"blind":"bigBlind","player":{"blind":"","cards":[],"intent":null,"lastGameAction":"","lastGameBet":0,"picture":"https://storage.googleapis.com/pokerblow-avatars/DIzWERaQmnZUxq9OBq4QXGigdpy2","position":1,"stack":296,"status":"playing","totalRoundBet":0,"userId":"DIzWERaQmnZUxq9OBq4QXGigdpy2","username":"Wonder_Woman"},"position":1}],"status":"playing","tournamentAttrs":{"levelIncreaseAt":1622462924735,"lobbyId":"60b4d10c039f12529aaeabfd","name":"Mistral","nextSmallBlind":2,"prizes":[{"place":1,"prize":600}],"startAt":1622462744735},"winners":[]}}},{"type":"blinds","payload":{"table":{"bettingLimitChips":9223372036854776000,"bigBlind":2,"maxRoundBet":2,"seats":[{"blind":"bigBlind","player":{"bankBalance":276161,"bankRank":2,"blind":"bigBlind","intent":null,"lastGameAction":"","lastGameBet":2,"level":6,"marketItemId":"lightBulb","picture":"https://storage.googleapis.com/pokerblow-avatars/DIzWERaQmnZUxq9OBq4QXGigdpy2","position":1,"stack":294,"status":"playing","totalRoundBet":2,"userId":"DIzWERaQmnZUxq9OBq4QXGigdpy2","username":"Wonder_Woman"},"position":1},{"blind":"dealerSmallBlind","player":{"bankBalance":266365,"bankRank":3,"blind":"dealerSmallBlind","intent":null,"lastGameAction":"","lastGameBet":1,"level":6,"marketItemId":"pacifier","picture":"https://storage.googleapis.com/pokerblow-avatars/cLhhrCC0sMUdpIHGlngodn1lKLl2","position":0,"stack":303,"status":"playing","totalRoundBet":1,"userId":"cLhhrCC0sMUdpIHGlngodn1lKLl2","username":"Sea_Man"},"position":0}],"smallBlind":1,"status":"playing","totalPot":3,"tournamentAttrs":{"levelIncreaseAt":1622462924735,"lobbyId":"60b4d10c039f12529aaeabfd","name":"Mistral","nextSmallBlind":2,"prizes":[{"place":1,"prize":600}],"startAt":1622462744735}}}},{"type":"holeCards","payload":{"table":{"seats":[{"blind":"dealerSmallBlind","player":{"bankBalance":266365,"bankRank":3,"blind":"dealerSmallBlind","cards":["Xx","Xx"],"intent":null,"lastGameAction":"","lastGameBet":1,"level":6,"marketItemId":"pacifier","picture":"https://storage.googleapis.com/pokerblow-avatars/cLhhrCC0sMUdpIHGlngodn1lKLl2","position":0,"stack":303,"status":"playing","totalRoundBet":1,"userId":"cLhhrCC0sMUdpIHGlngodn1lKLl2","username":"Sea_Man"},"position":0},{"blind":"bigBlind","player":{"bankBalance":276161,"bankRank":2,"blind":"bigBlind","cards":["7d","6c"],"intent":null,"lastGameAction":"","lastGameBet":2,"level":6,"marketItemId":"lightBulb","picture":"https://storage.googleapis.com/pokerblow-avatars/DIzWERaQmnZUxq9OBq4QXGigdpy2","position":1,"stack":294,"status":"playing","totalRoundBet":2,"userId":"DIzWERaQmnZUxq9OBq4QXGigdpy2","username":"Wonder_Woman"},"position":1}]}}},{"type":"timeToDecide","payload":{"table":{"decidingPosition":0,"lastAggressorPosition":-1,"seats":[{"blind":"dealerSmallBlind","player":{"intent":null,"isDeciding":true,"position":0,"timeoutAt":1622462846769},"position":0}]}}}]`
	merge(resetHoleCards)
	userAction := `[{"type":"playerAction","payload":{"table":{"bettingLimitChips":9223372036854776000,"decidingPosition":-1,"maxRoundBet":90,"seats":[{"blind":"dealerSmallBlind","player":{"bankBalance":266365,"bankRank":3,"blind":"dealerSmallBlind","intent":null,"isDeciding":false,"lastGameAction":"raise","lastGameBet":90,"level":6,"marketItemId":"pacifier","picture":"https://storage.googleapis.com/pokerblow-avatars/cLhhrCC0sMUdpIHGlngodn1lKLl2","position":0,"stack":214,"status":"playing","totalRoundBet":90,"userId":"cLhhrCC0sMUdpIHGlngodn1lKLl2","username":"Sea_Man"},"position":0}],"totalPot":92}}},{"type":"timeToDecide","payload":{"table":{"decidingPosition":1,"lastAggressorPosition":0,"seats":[{"blind":"bigBlind","player":{"intent":null,"isDeciding":true,"position":1,"timeoutAt":1622462851129},"position":1}]}}}]`
	merge(userAction)
	botAction := `[{"type":"playerAction","payload":{"table":{"bettingLimitChips":9223372036854776000,"decidingPosition":-1,"maxRoundBet":0,"seats":[{"blind":"bigBlind","player":{"bankBalance":276161,"bankRank":2,"blind":"bigBlind","intent":null,"isDeciding":false,"lastGameAction":"call","lastGameBet":90,"level":6,"marketItemId":"lightBulb","picture":"https://storage.googleapis.com/pokerblow-avatars/DIzWERaQmnZUxq9OBq4QXGigdpy2","position":1,"stack":206,"status":"playing","totalRoundBet":90,"userId":"DIzWERaQmnZUxq9OBq4QXGigdpy2","username":"Wonder_Woman"},"position":1}],"totalPot":180}}},{"type":"newBettingRound","payload":{"newCards":["2d","9c","8c"],"roundType":"flop","table":{"communityCards":["2d","9c","8c"],"pots":[{"chips":180,"winnerPositions":null}],"seats":[{"blind":"dealerSmallBlind","player":{"intent":null,"position":0,"totalRoundBet":0},"position":0},{"blind":"bigBlind","player":{"intent":null,"position":1,"totalRoundBet":0},"position":1}],"totalPot":180}}},{"type":"timeToDecide","payload":{"table":{"decidingPosition":1,"lastAggressorPosition":-1,"seats":[{"blind":"bigBlind","player":{"intent":null,"isDeciding":true,"position":1,"timeoutAt":1622462853761},"position":1}]}}}]`
	merge(botAction)
	botActionFlop := `[{"type":"playerAction","payload":{"table":{"bettingLimitChips":9223372036854776000,"decidingPosition":-1,"maxRoundBet":0,"seats":[{"blind":"bigBlind","player":{"bankBalance":276161,"bankRank":2,"blind":"bigBlind","intent":null,"isDeciding":false,"lastGameAction":"check","lastGameBet":0,"level":6,"marketItemId":"lightBulb","picture":"https://storage.googleapis.com/pokerblow-avatars/DIzWERaQmnZUxq9OBq4QXGigdpy2","position":1,"stack":206,"status":"playing","totalRoundBet":0,"userId":"DIzWERaQmnZUxq9OBq4QXGigdpy2","username":"Wonder_Woman"},"position":1}],"totalPot":180}}},{"type":"timeToDecide","payload":{"table":{"decidingPosition":0,"lastAggressorPosition":-1,"seats":[{"blind":"dealerSmallBlind","player":{"intent":null,"isDeciding":true,"position":0,"timeoutAt":1622462855435},"position":0}]}}}]`
	merge(botActionFlop)
	userActionFlop := `[{"type":"playerAction","payload":{"table":{"bettingLimitChips":9223372036854776000,"decidingPosition":-1,"maxRoundBet":0,"seats":[{"blind":"dealerSmallBlind","player":{"bankBalance":266365,"bankRank":3,"blind":"dealerSmallBlind","intent":null,"isDeciding":false,"lastGameAction":"check","lastGameBet":0,"level":6,"marketItemId":"pacifier","picture":"https://storage.googleapis.com/pokerblow-avatars/cLhhrCC0sMUdpIHGlngodn1lKLl2","position":0,"stack":214,"status":"playing","totalRoundBet":0,"userId":"cLhhrCC0sMUdpIHGlngodn1lKLl2","username":"Sea_Man"},"position":0}],"totalPot":180}}},{"type":"newBettingRound","payload":{"newCards":["5s"],"roundType":"turn","table":{"communityCards":["2d","9c","8c","5s"],"pots":[{"chips":180,"winnerPositions":null}],"seats":[{"blind":"dealerSmallBlind","player":{"intent":null,"position":0,"totalRoundBet":0},"position":0},{"blind":"bigBlind","player":{"intent":null,"position":1,"totalRoundBet":0},"position":1}],"totalPot":180}}},{"type":"timeToDecide","payload":{"table":{"decidingPosition":1,"lastAggressorPosition":-1,"seats":[{"blind":"bigBlind","player":{"intent":null,"isDeciding":true,"position":1,"timeoutAt":1622462861701},"position":1}]}}}]`
	merge(userActionFlop)
	action := (&Bot{Position: 1, Looseness: 0.5, Aggression: 0.8, prevAction: Check}).GetAction(table)
	assert.EqualValues(t, AllInType, action.Type)
}

func extractEvents(data string) (events []*Event) {
	for _, result := range gjson.Get(data, "@this").Array() {
		events = append(events, NewEventBytes([]byte(result.String())))
	}
	return
}
