package domain

import (
	"encoding/json"
	"github.com/glossd/pokergloss/survival/bot-squad/conf"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

var defaultConfig = conf.Config{
	TableID:      "60f97044fd07b00f284e515a",
	Token:        "eyJhbGciOiJSUzI1NiIsImtpZCI6IjFiYjk2MDVjMzZlOThlMzAxMTdhNjk1MTc1NjkzODY4MzAyMDJiMmQiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiUG9rZXJibG93IiwicGljdHVyZSI6Imh0dHBzOi8vc3RvcmFnZS5nb29nbGVhcGlzLmNvbS9wb2tlcmJsb3ctYXZhdGFycy9TOWRKM0hIUHNUWHhScWx6SUFiQTdieHgzTWEyLUtQZFkiLCJ1c2VybmFtZSI6IlBva2VyYmxvdyIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS9wb2tlcmJsb3ciLCJhdWQiOiJwb2tlcmJsb3ciLCJhdXRoX3RpbWUiOjE2MjI5MTEwMDAsInVzZXJfaWQiOiJTOWRKM0hIUHNUWHhScWx6SUFiQTdieHgzTWEyIiwic3ViIjoiUzlkSjNISFBzVFh4UnFseklBYkE3Ynh4M01hMiIsImlhdCI6MTYyNjk1OTg0OCwiZXhwIjoxNjI2OTYzNDQ4LCJlbWFpbCI6InBva2VyYmxvd0Bwb2tlcmJsb3cuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImZpcmViYXNlIjp7ImlkZW50aXRpZXMiOnsiZW1haWwiOlsicG9rZXJibG93QHBva2VyYmxvdy5jb20iXX0sInNpZ25faW5fcHJvdmlkZXIiOiJwYXNzd29yZCJ9fQ.Pf886VzqiBUHHJcrpG9CsBiXgAf6_Q4cwDgXXB-gCZJ7SQncfGb76gNVHISdkVuyRI4m61KE0E2W4nALCnCDU6_v3Uah3q-nDteb7NYk1o1JoWTC_bs_LbO2STgKQOXvIU6YPcx8k-b7--vsOmH4TRRPN-STrlcHf-B0ao0fVl6kyrJ4Yh7t4M48QgfPLZWuMzPbNbkBUHlDbzFMmQwa70H1Vg_-Fq0FZiz5Xqsmyf-oIM-_np2AGPjIvAeQ5IHDQTm54qIJDvNpJAX90GLtk1wynetut0HzHbeAjs6-kjrL0I0qvzZWn62uBvVErz6CYu1-mK7byimO30E-mdmVKw",
	UserPosition: 0,
	Squad: conf.Squad{
		Positions:  []string{"1"},
		Looseness:  []string{"0.5"},
		Aggression: []string{"0.75"},
	}}

func TestNewEvent(t *testing.T) {
	e := NewEventBytes(timeToDecide())
	assert.EqualValues(t, "timeToDecide", e.Type)
	assert.NotZero(t, len(e.Payload))
}

func TestMergeJSON(t *testing.T) {
	table := getFullTable()
	assert.EqualValues(t, 3, table.TotalPot)
	assert.EqualValues(t, 1, table.Seats[0].Player.TotalRoundBet)
	assert.EqualValues(t, 0, table.DecidingPosition)
	assert.EqualValues(t, 2, len(table.Seats[0].Player.HoleCards))
	tmerge(t, table, `{"type":"playerAction","payload":{"table":{"bettingLimitChips":9223372036854776000,"decidingPosition":-1,"maxRoundBet":2,"seats":[{"blind":"dealerSmallBlind","player":{"bankBalance":96882,"bankRank":12,"blind":"dealerSmallBlind","intent":null,"isDeciding":false,"lastGameAction":"call","lastGameBet":2,"level":6,"picture":"https://storage.googleapis.com/pokerblow-avatars/cLhhrCC0sMUdpIHGlngodn1lKLl2","position":0,"stack":249,"status":"playing","totalRoundBet":2,"userId":"cLhhrCC0sMUdpIHGlngodn1lKLl2","username":"Sea_Man"},"position":0}],"totalPot":4}}}`)
	assert.EqualValues(t, 4, table.TotalPot)
	assert.EqualValues(t, 2, table.Seats[0].Player.TotalRoundBet)
	tmerge(t, table, `{"type":"timeToDecide","payload":{"table":{"decidingPosition":1,"lastAggressorPosition":-1,"seats":[{"blind":"bigBlind","player":{"intent":null,"isDeciding":true,"position":1,"timeoutAt":1620164737289},"position":1}]}}`)
	assert.EqualValues(t, 1, table.DecidingPosition)
	assert.EqualValues(t, 2, table.MaxRoundBet)
	tmerge(t, table, `{"type":"playerAction","payload":{"table":{"bettingLimitChips":9223372036854776000,"decidingPosition":-1,"maxRoundBet":0,"seats":[{"blind":"bigBlind","player":{"bankBalance":1950,"bankRank":18,"blind":"bigBlind","intent":null,"isDeciding":false,"lastGameAction":"check","lastGameBet":2,"level":1,"picture":"https://storage.googleapis.com/pokerblow-avatars/ART2wmVl6MblilWIQIJkWG6ztWW2","position":1,"stack":248,"status":"playing","totalRoundBet":2,"userId":"ART2wmVl6MblilWIQIJkWG6ztWW2","username":"pokerblow"},"position":1}],"totalPot":4}}}`)
	assert.EqualValues(t, 0, len(table.CommCards))
	tmerge(t, table, `{"type":"newBettingRound","payload":{"newCards":["Jh","Ks","Qs"],"roundType":"flop","table":{"communityCards":["Jh","Ks","Qs"],"pots":[{"chips":4,"winnerPositions":null}],"seats":[{"blind":"dealerSmallBlind","player":{"intent":null,"position":0,"totalRoundBet":0},"position":0},{"blind":"bigBlind","player":{"intent":null,"position":1,"totalRoundBet":0},"position":1}],"totalPot":4}}}`)
	assert.EqualValues(t, 3, len(table.CommCards))
	assert.EqualValues(t, 0, table.MaxRoundBet)
	tmerge(t, table, `{"type":"timeToDecide","payload":{"table":{"decidingPosition":1,"lastAggressorPosition":-1,"seats":[{"blind":"bigBlind","player":{"intent":null,"isDeciding":true,"position":1,"timeoutAt":1620166019554},"position":1}]}}}`)

	tmerge(t, table, `{"type":"newBettingRound","payload":{"newCards":["5s"],"roundType":"turn","table":{"communityCards":["Td","Qs","Ks","5s"],"pots":[{"chips":4,"winnerPositions":null}],"seats":[{"blind":"dealerSmallBlind","player":{"intent":null,"position":0,"totalRoundBet":0},"position":0},{"blind":"bigBlind","player":{"intent":null,"position":1,"totalRoundBet":0},"position":1}],"totalPot":4}}}`)
	assert.EqualValues(t, 4, len(table.CommCards))

	tmerge(t, table, `{"type":"reset","payload":{"table":{"bettingLimitChips":0,"communityCards":[],"decidingPosition":-1,"lastAggressorPosition":-1,"maxRoundBet":0,"pots":[],"seats":[{"blind":"","player":{"blind":"","cards":[],"intent":null,"lastGameAction":"","lastGameBet":0,"picture":"https://storage.googleapis.com/pokerblow-avatars/cLhhrCC0sMUdpIHGlngodn1lKLl2","position":0,"stack":252,"status":"ready","totalRoundBet":0,"userId":"cLhhrCC0sMUdpIHGlngodn1lKLl2","username":"Sea_Man"},"position":0},{"blind":"","player":{"blind":"","cards":[],"intent":null,"lastGameAction":"","lastGameBet":0,"picture":"https://storage.googleapis.com/pokerblow-avatars/ART2wmVl6MblilWIQIJkWG6ztWW2","position":1,"stack":248,"status":"sittingOut","totalRoundBet":0,"userId":"ART2wmVl6MblilWIQIJkWG6ztWW2","username":"pokerblow"},"position":1}],"status":"waiting","winners":[]}}}`)
	assert.EqualValues(t, 0, len(table.CommCards))
	tmerge(t, table, `{"type":"holeCards","payload":{"table":{"seats":[{"blind":"bigBlind","player":{"bankBalance":96905,"bankRank":12,"blind":"bigBlind","cards":["Xx","Xx"],"intent":null,"lastGameAction":"","lastGameBet":2,"level":6,"picture":"https://storage.googleapis.com/pokerblow-avatars/cLhhrCC0sMUdpIHGlngodn1lKLl2","position":0,"stack":250,"status":"playing","totalRoundBet":2,"userId":"cLhhrCC0sMUdpIHGlngodn1lKLl2","username":"Sea_Man"},"position":0},{"blind":"dealerSmallBlind","player":{"blind":"dealerSmallBlind","cards":["Ac","Kh"],"intent":null,"lastGameAction":"","lastGameBet":1,"picture":"https://storage.googleapis.com/pokerblow-avatars/ART2wmVl6MblilWIQIJkWG6ztWW2","position":1,"stack":249,"status":"playing","totalRoundBet":1,"userId":"ART2wmVl6MblilWIQIJkWG6ztWW2","username":"pokerblow"},"position":1}]}}}`)
	assert.EqualValues(t, "Ac", table.Seats[1].Player.HoleCards[0].String())

	tmerge(t, table, `{"type":"playerLeft", "payload":{"table":{"seats":[{"position":0, "blind":"", "player":null}]}}}`)
	assert.Nil(t, table.Seats[0].Player)
}

func tmerge(t *testing.T, table *Table, event string) {
	assert.Nil(t, NewEventBytes([]byte(event)).Merge(defaultConfig, table))
}

func TestMergeTable(t *testing.T) {
	table := getEmptyTable(t)
	assert.Nil(t, table.Seats[0].Player)

	assert.Nil(t, MergeTableBytes(defaultConfig, table, seatReserved()))
	assert.NotNil(t, table.Seats[0].Player)
	assert.EqualValues(t, firstUserId, table.Seats[0].Player.UserID)
	assert.EqualValues(t, 6, len(table.Seats))
}

func timeToDecide() []byte {
	return []byte(`{"type":"timeToDecide","payload":{"table":{"decidingPosition":1,"lastAggressorPosition":-1,"seats":[{"blind":"bigBlind","player":{"intent":null,"isDeciding":true,"position":1,"timeoutAt":1620164737289},"position":1}]}}`)
}

func seatReserved() []byte {
	return []byte(`{"type":"seatReserved","payload":{"table":{"seats":[{"blind":"","player":{"blind":"","intent":null,"lastGameAction":"","lastGameBet":0,"picture":"https://storage.googleapis.com/pokerblow-avatars/DIzWERaQmnZUxq9OBq4QXGigdpy2","position":0,"stack":0,"status":"sitting","totalRoundBet":0,"userId":"DIzWERaQmnZUxq9OBq4QXGigdpy2","username":"Wonder_Woman"},"position":0}]}}}`)
}

func getFullTable() *Table {
	var t Table
	err := json.Unmarshal(fullTable(), &t)
	if err != nil {
		log.Panicf("Failed to unmarshal the table: %s", err)
	}
	return &t
}

func fullTable() []byte {
	return []byte(`{"id":"6091befdc1822a5193a60427","name":"Chinook","size":2,"bigBlind":2,"smallBlind":1,"decisionTimeoutSec":10,"bettingLimit":"NL","minBuyIn":100,"maxBuyIn":400,"type":"cashGame","occupied":2,"avgStake":0,"avgPot":0,"seats":[{"position":0,"blind":"dealerSmallBlind","player":{"userId":"cLhhrCC0sMUdpIHGlngodn1lKLl2","username":"Sea_Man","picture":"https://storage.googleapis.com/pokerblow-avatars/cLhhrCC0sMUdpIHGlngodn1lKLl2","position":0,"stack":250,"status":"playing","cards":["Xx","Xx"],"blind":"dealerSmallBlind","totalRoundBet":1,"lastGameBet":1,"lastGameAction":"","timeoutAt":1620164395101,"intent":null,"level":6,"bankBalance":96882,"bankRank":12,"isDeciding":true}},{"position":1,"blind":"bigBlind","player":{"userId":"ART2wmVl6MblilWIQIJkWG6ztWW2","username":"pokerblow","picture":"https://storage.googleapis.com/pokerblow-avatars/ART2wmVl6MblilWIQIJkWG6ztWW2","position":1,"stack":247,"status":"playing","cards":["6d","4c"],"blind":"bigBlind","totalRoundBet":2,"lastGameBet":2,"lastGameAction":"","intent":null,"level":1,"bankBalance":971,"bankRank":19}}],"status":"playing","maxRoundBet":2,"bettingLimitChips":9223372036854775807,"pots":null,"totalPot":3,"communityCards":[],"decidingPosition":0,"lastAggressorPosition":-1}`)
}
