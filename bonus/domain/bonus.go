package domain

import (
	"math"
)

const firstDayBonus float64 = 500.0

type DailyBonus struct {
	UserID string `bson:"_id"`
	DayInARow int
	IsTaken bool
}

func NewDailyBonus(userID string) *DailyBonus {
	return &DailyBonus{
		UserID:    userID,
		DayInARow: 0,
		IsTaken:   false,
	}
}

func (b *DailyBonus) Visit() bool {
	if !b.IsTaken {
		b.DayInARow++
		b.IsTaken = true
		return true
	}
	return false
}

func (b *DailyBonus) IsVisited() bool {
	return b.IsTaken
}

func (b *DailyBonus) Reset() {
	if !b.IsTaken {
		b.DayInARow = 0
	}
	b.IsTaken = false
}

// Valid after Visit
func (b *DailyBonus) CalculateBonus() int64 {
	if b.DayInARow == 0 {
		return int64(firstDayBonus)
	}
	return int64(firstDayBonus + math.Sqrt(float64(2 * (b.DayInARow - 1)) ) * firstDayBonus)
}


