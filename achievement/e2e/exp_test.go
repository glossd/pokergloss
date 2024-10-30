package e2e

import (
	"github.com/glossd/pokergloss/achievement/db"
	"github.com/glossd/pokergloss/achievement/service"
	"github.com/glossd/pokergloss/gomq/mqtable"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"testing"
)

func TestFirstExp(t *testing.T) {
	t.Cleanup(cleanUp)
	gameEnd := oneWinnerTwoPlayers()
	service.UpdateExp(gameEnd)
	assertPoints(t, defaultIdentity.UserId, 2)
	assertPoints(t, "2", 1)

	rr := testRouter.Request(t, http.MethodGet, "/points/me", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	assert.EqualValues(t, 2, gjson.Get(body, "points").Int())
	assert.EqualValues(t, 1, gjson.Get(body, "level").Int())
}

func oneWinnerTwoPlayers() *mqtable.GameEnd {
	return &mqtable.GameEnd{
		Winners:        []*mqtable.Winner{{UserId: defaultIdentity.UserId, Chips: 4, Hand: "Straight"}},
		Players:        []*mqtable.Player{{UserId: defaultIdentity.UserId, WageredChips: 2, LastAction: "allIn"}, {UserId: "2", WageredChips: 2}},
		CommunityCards: []string{"As", "Kd", "Qs", "Td", "7s"},
	}
}

func oneDefaultLosesTwoPlayers() *mqtable.GameEnd {
	return &mqtable.GameEnd{
		Winners:        []*mqtable.Winner{{UserId: "2", Chips: 4, Hand: "Straight"}},
		Players:        []*mqtable.Player{{UserId: defaultIdentity.UserId, WageredChips: 2, LastAction: "fold"}, {UserId: "2", WageredChips: 2}},
		CommunityCards: []string{"As", "Kd", "Qs", "Td", "7s"},
	}
}

func oneWinnerPlayers(players []*mqtable.Player) *mqtable.GameEnd {
	return &mqtable.GameEnd{
		Winners:        []*mqtable.Winner{{UserId: defaultIdentity.UserId, Chips: 4, Hand: "Straight"}},
		Players:        players,
		CommunityCards: []string{"As", "Kd", "Qs", "Td", "7s"},
	}
}

func assertPoints(t *testing.T, usereID string, points int64) {
	exp, err := db.FindExpNoCtx(usereID)
	assert.Nil(t, err)
	assert.EqualValues(t, points, exp.Points)
}
