package domain

type Counter interface {
	GetPrize() Prize
	GetType() AchievementType
	GetName() string
	GetCount() int64
	GetLevel() int
	LevelCount() int64
	NextLevelCount() int64
	NextLevelPrize() int64
}
