package domain

import (
	"github.com/glossd/pokergloss/gomq/mqtable"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBustCounter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dc := NewBustCounter()
		winner := &mqtable.Winner{UserId: "1"}
		ge := &mqtable.GameEnd{
			Winners: []*mqtable.Winner{winner},
			Players: []*mqtable.Player{pAllIn("2"), pAllIn("3")}}
		dc.Update(winner, NewGameEnd(ge))
		assert.EqualValues(t, 1, dc.Count)
		assert.NotZero(t, dc.GetPrize().Chips)
	})
	t.Run("skip", func(t *testing.T) {
		dc := NewDefeatCounter()
		winner := &mqtable.Winner{UserId: "1"}
		ge := &mqtable.GameEnd{
			Winners: []*mqtable.Winner{winner},
			Players: []*mqtable.Player{pAllIn("2"), {UserId: "3", LastAction: "fold"}}}
		dc.Update(winner, NewGameEnd(ge))
		assert.EqualValues(t, 0, dc.Count)
	})
}

func pAllIn(userId string) *mqtable.Player {
	return &mqtable.Player{UserId: userId, LastAction: "allIn", WageredChips: 234}
}
