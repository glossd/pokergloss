package domain

import (
	"github.com/glossd/pokergloss/gomq/mqtable"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComputeLevel(t *testing.T) {
	assert.EqualValues(t, 1, computeLevel(0))
	assert.EqualValues(t, 1, computeLevel(1))
	assert.EqualValues(t, 1, computeLevel(5))
	assert.EqualValues(t, 1, computeLevel(10))
	assert.EqualValues(t, 1, computeLevel(19))
	assert.EqualValues(t, 2, computeLevel(20))
	assert.EqualValues(t, 2, computeLevel(21))
	assert.EqualValues(t, 2, computeLevel(29))
	assert.EqualValues(t, 3, computeLevel(40))
	assert.EqualValues(t, 3, computeLevel(40))
}

func TestComputeStartLevelPoints(t *testing.T) {
	assert.EqualValues(t, 0, computeStartLevelPoints(1))
	assert.EqualValues(t, 20, computeStartLevelPoints(2))
	assert.EqualValues(t, 40, computeStartLevelPoints(3))
	assert.EqualValues(t, 80, computeStartLevelPoints(4))
	assert.EqualValues(t, 160, computeStartLevelPoints(5))
}

func TestComputeLevelPrize(t *testing.T) {
	assert.EqualValues(t, 1000, computeExPLevelPrize(2))
	assert.EqualValues(t, 2000, computeExPLevelPrize(3))
	assert.EqualValues(t, 4000, computeExPLevelPrize(4))
	assert.EqualValues(t, 8000, computeExPLevelPrize(5))
}

func TestComputeExP(t *testing.T) {
	gameEnd := &mqtable.GameEnd{
		Winners:        []*mqtable.Winner{{UserId: "1", Chips: 4, Hand: "Straight"}},
		Players:        []*mqtable.Player{{UserId: "1", WageredChips: 2}, {UserId: "2", WageredChips: 2}},
		CommunityCards: []string{"As", "Kd", "Qs", "Td", "7s"},
	}

	exp := NewExP("1")
	exp.UpdateWithGameEnd(gameEnd)
	assert.EqualValues(t, 2, exp.Points)
	assert.EqualValues(t, 1, exp.Level)
}
