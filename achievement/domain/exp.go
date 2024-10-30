package domain

import (
	"github.com/glossd/pokergloss/gomq/mqtable"
	"math"
)

// Experience points
type ExP struct {
	UserID string `bson:"_id"`
	Points int64
	Level  int

	isNewLevel bool
}

func NewExP(userID string) *ExP {
	return &ExP{
		UserID: userID,
		Points: 0,
		Level:  1,
	}
}

func (e *ExP) UpdateWithGameEnd(end *mqtable.GameEnd) {
	e.Increase(e.ComputePoints(NewGameEnd(end)))
}

func (e *ExP) Increase(points int64) {
	e.Points += points
	oldLevel := e.Level
	e.Level = computeLevel(e.Points)
	e.isNewLevel = oldLevel != e.Level
}

func (e *ExP) ComputePoints(end *GameEnd) int64 {
	if len(end.CommunityCards) < 3 {
		return 0
	}
	p, ok := end.PlayersMap[e.UserID]
	if !ok {
		return 0
	}
	if p.WageredChips == 0 {
		return 0
	}

	_, ok = end.WinnersMap[e.UserID]
	if ok {
		return int64(len(end.WageredPlayers))
	}
	return 1
}

func (e *ExP) IsNewLevel() bool {
	return e.isNewLevel
}

func (e *ExP) NextLevelPoints() int64 {
	return computeStartLevelPoints(e.Level + 1)
}

func (e *ExP) StartLevelPoints() int64 {
	return computeStartLevelPoints(e.Level)
}

func computeLevel(points int64) int {
	if points < 20 {
		return 1
	}
	return int(math.Floor(math.Log2(float64(points)/10))) + 1
}

func computeStartLevelPoints(level int) int64 {
	if level == 1 {
		return 0
	}
	return int64(math.Pow(2, float64(level-1)) * 10)
}

func (e *ExP) GetNewLevelPrize() int64 {
	if e.isNewLevel {
		return computeExPLevelPrize(e.Level)
	}
	return 0
}

func computeExPLevelPrize(level int) int64 {
	if level < 2 {
		return 0
	}
	return int64(math.Pow(2, float64(level-2))) * 1000
}
