package domain

import (
	"github.com/glossd/pokergloss/gomq/mqtable"
)

type AchievementStore struct {
	UserID           string `bson:"_id"`
	HandsCounter     *HandsCounter
	WinCounter       *WinCounter
	SitngoWinCounter *SitngoWinCounter
	MultiWinCounter  *MultiWinCounter
	BustCounter      *BustCounter
	DefeatCounter    *DefeatCounter
}

func NewAchievementStore(userID string) *AchievementStore {
	return &AchievementStore{
		UserID:           userID,
		HandsCounter:     NewHandsCounter(),
		WinCounter:       NewWinCounter(),
		SitngoWinCounter: NewSitngoWinCounter(),
		MultiWinCounter:  NewMultiWinCounter(),
		DefeatCounter:    NewDefeatCounter(),
		BustCounter:      NewBustCounter(),
	}
}

func (as *AchievementStore) Update(winner *mqtable.Winner, ge *GameEnd) {
	_ = as.HandsCounter.Update(winner)
	as.WinCounter.Inc()
	if as.BustCounter == nil {
		// for the time between states deployed and ran migrations
		as.BustCounter = NewBustCounter()
	}
	as.BustCounter.Update(winner, ge)
	if as.DefeatCounter == nil {
		as.DefeatCounter = NewDefeatCounter()
	}
}

func (as *AchievementStore) GetCounters() []Counter {
	var a []Counter
	if as.WinCounter != nil {
		a = append(a, as.WinCounter)
	}
	if as.SitngoWinCounter != nil {
		a = append(a, as.SitngoWinCounter)
	}
	if as.MultiWinCounter != nil {
		a = append(a, as.MultiWinCounter)
	}
	if as.DefeatCounter != nil {
		a = append(a, as.DefeatCounter)
	}
	if as.BustCounter != nil {
		a = append(a, as.BustCounter)
	}
	return a
}
