package domain

import (
	"github.com/glossd/pokergloss/gomq/mqtable"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefeatCounter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dc := NewDefeatCounter()
		winner := &mqtable.Winner{UserId: "1"}
		ge := &mqtable.GameEnd{
			Winners: []*mqtable.Winner{winner},
			Players: []*mqtable.Player{pCall("2"), pCall("3"), pCall("4")}}
		dc.Update(winner, NewGameEnd(ge))
		assert.EqualValues(t, 1, dc.Count)
		assert.NotZero(t, dc.GetPrize().Chips)
	})
	t.Run("skip", func(t *testing.T) {
		dc := NewDefeatCounter()
		winner := &mqtable.Winner{UserId: "1"}
		ge := &mqtable.GameEnd{
			Winners: []*mqtable.Winner{winner},
			Players: []*mqtable.Player{pCall("2"), pCall("3")}}
		dc.Update(winner, NewGameEnd(ge))
		assert.EqualValues(t, 0, dc.Count)
	})
	t.Run("skip/foldedPlayers", func(t *testing.T) {
		dc := NewDefeatCounter()
		winner := &mqtable.Winner{UserId: "1"}
		ge := &mqtable.GameEnd{
			Winners: []*mqtable.Winner{winner},
			Players: []*mqtable.Player{{UserId: "2", LastAction: "fold"}, pCall("3"), pCall("4")}}
		dc.Update(winner, NewGameEnd(ge))
		assert.EqualValues(t, 0, dc.Count)
	})
}

func pCall(userId string) *mqtable.Player {
	return &mqtable.Player{UserId: userId, LastAction: "call", WageredChips: 234}
}
