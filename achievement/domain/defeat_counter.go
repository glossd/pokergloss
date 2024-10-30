package domain

import "github.com/glossd/pokergloss/gomq/mqtable"

type DefeatCounter struct {
	SimpleCounter `bson:",inline"`
}

func (b *DefeatCounter) GetType() AchievementType { return SimpleCounterType }
func (b *DefeatCounter) GetName() string          { return "Defeat 3+ Players" }

func NewDefeatCounter() *DefeatCounter {
	return &DefeatCounter{SimpleCounter: NewSimpleCounter(50, 5, true)}
}

func (b *DefeatCounter) Update(w *mqtable.Winner, ge *GameEnd) {
	if b == nil {
		return
	}
	if len(ge.Winners) != 1 {
		return
	}
	if w.UserId != ge.Winners[0].UserId {
		return
	}
	lostPlayers := ge.TillTheEndPlayersExceptWinners()
	if len(lostPlayers) < 3 {
		return
	}
	b.Inc()
}
