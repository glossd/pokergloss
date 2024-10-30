package domain

type AchievementType string
const (
	HandType AchievementType = "hand"
	SimpleCounterType AchievementType = "simpleCounter"
)


// Maybe one day
type Achievement interface {
	LevelCount() int64
	NextLevelCount() int64
	NextLevelPrize() int64
	GetPrize() Prize
	GetLevel() int
	GetCount() int64

	countByLevel(lvl int) int64
	updateLevel()
	setPrize()

	GetName() string
}
