package domain

import "github.com/glossd/pokergloss/gomq/mqtable"

type BustCounter struct {
	SimpleCounter `bson:",inline"`
}

func (b *BustCounter) GetType() AchievementType { return SimpleCounterType }
func (b *BustCounter) GetName() string          { return "Bust 2+ Players" }

func NewBustCounter() *BustCounter {
	return &BustCounter{SimpleCounter: NewSimpleCounter(50, 5, true)}
}

func (b *BustCounter) Update(w *mqtable.Winner, ge *GameEnd) {
	if b == nil {
		return
	}
	if len(ge.Winners) != 1 {
		return
	}
	if w.UserId != ge.Winners[0].UserId {
		return
	}
	allInPlayers := ge.AllInWageredPlayersExceptWinners()
	if len(allInPlayers) < 2 {
		return
	}
	b.Inc()
}
