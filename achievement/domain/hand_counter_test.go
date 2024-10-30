package domain

import (
	"github.com/glossd/pokergloss/gomq/mqtable"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandsCounter_Update(t *testing.T) {
	hc := NewHandsCounter()
	err := hc.Update(&mqtable.Winner{
		Chips: 4,
		Hand:  "Straight",
	})
	assert.Nil(t, err)

	assert.EqualValues(t, 7, len(hc.Hands))
	assert.EqualValues(t, 1, hc.Hands[S].Count)
	assert.EqualValues(t, 500, hc.GetPrize().Chips)
}

func TestHandsCounter_Update_NoneCountedHand(t *testing.T) {
	hc := NewHandsCounter()
	err := hc.Update(&mqtable.Winner{
		Chips: 4,
		Hand:  "High Card",
	})
	assert.Nil(t, err)

	assert.EqualValues(t, 7, len(hc.Hands))
	assert.EqualValues(t, 0, hc.GetPrize().Chips)
}

func TestHandsCounter_Update_UnknownHand(t *testing.T) {
	hc := NewHandsCounter()
	err := hc.Update(&mqtable.Winner{
		Chips: 4,
		Hand:  "Blah",
	})
	assert.NotNil(t, err)
}

func TestComputeHandAchievementLevel(t *testing.T) {
	assert.EqualValues(t, 0, computeHandAchievementLevel(0))
	assert.EqualValues(t, 1, computeHandAchievementLevel(1))
	assert.EqualValues(t, 1, computeHandAchievementLevel(2))
	assert.EqualValues(t, 1, computeHandAchievementLevel(4))
	assert.EqualValues(t, 2, computeHandAchievementLevel(5))
	assert.EqualValues(t, 2, computeHandAchievementLevel(24))
	assert.EqualValues(t, 3, computeHandAchievementLevel(25))
}

func TestGetLevelCount(t *testing.T) {
	assert.EqualValues(t, 0, GetLevelCount(0))
	assert.EqualValues(t, 1, GetLevelCount(1))
	assert.EqualValues(t, 5, GetLevelCount(2))
	assert.EqualValues(t, 25, GetLevelCount(3))
	assert.EqualValues(t, 125, GetLevelCount(4))
}

func TestHandsCounter_OrderedHands(t *testing.T) {
	hc := NewHandsCounter()
	hands := hc.OrderedCounters()
	assert.EqualValues(t, 7, len(hands))
	assert.EqualValues(t, TP, hands[0].Hand)
	assert.EqualValues(t, ToK, hands[1].Hand)
	assert.EqualValues(t, SF, hands[6].Hand)
}

func TestComputeHandAchievedLevelPrize(t *testing.T) {
	assert.EqualValues(t, 0, ComputeHandLevelPrize(0, S))
	assert.EqualValues(t, 500, ComputeHandLevelPrize(1, S))
	assert.EqualValues(t, 2500, ComputeHandLevelPrize(2, S))
}
