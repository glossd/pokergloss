package domain
import "math"

type SimpleCounter struct {
	Count int64
	Level int

	PrizePerUnit       int64
	LevelLogarithmBase int64
	// Defines whether counter increments level starting from the first inc
	FirstIncluded bool

	prize Prize
}

func NewSimpleCounter(prizePerUnit int64, levelLogarithmBase int64, firstIncluded bool) SimpleCounter {
	return SimpleCounter{PrizePerUnit: prizePerUnit, LevelLogarithmBase: levelLogarithmBase, FirstIncluded: firstIncluded}
}

func (sc *SimpleCounter) Inc() {
	sc.Count++
	oldLevel := sc.Level
	sc.updateLevel()
	if oldLevel != sc.Level {
		sc.setPrize()
	}
}

func (sc *SimpleCounter) GetPrize() Prize {
	return sc.prize
}

func (sc *SimpleCounter) updateLevel() {
	if sc.Count == 0 {
		sc.Level = 0
	}
	level := int(math.Floor(Logarithm(sc.Count, sc.LevelLogarithmBase)))
	if sc.FirstIncluded {
		sc.Level = level+1
	} else{
		sc.Level = level
	}
}

func (sc *SimpleCounter) setPrize() {
	sc.prize = Prize{
		Chips: sc.computeLevelPrize(sc.Level),
	}
}

func (sc *SimpleCounter) computeLevelPrize(level int) int64 {
	return sc.countByLevel(level) * sc.PrizePerUnit
}

func (sc *SimpleCounter) GetCount() int64 {
	return sc.Count
}

func (sc *SimpleCounter) GetLevel() int {
	return sc.Level
}

func (sc *SimpleCounter) LevelCount() int64 {
	return sc.countByLevel(sc.Level)
}

func (sc *SimpleCounter) NextLevelCount() int64 {
	return sc.countByLevel(sc.Level+1)
}

func (sc *SimpleCounter) NextLevelPrize() int64 {
	return sc.computeLevelPrize(sc.Level+1)
}

func (sc *SimpleCounter) countByLevel(lvl int) int64 {
	if lvl == 0 {
		return 0
	}
	level := lvl
	if sc.FirstIncluded {
		level = lvl-1
	}
	return int64(math.Pow(float64(sc.LevelLogarithmBase), float64(level)))
}
