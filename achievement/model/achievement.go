package model

import "github.com/glossd/pokergloss/achievement/domain"

type Achievement struct {
	Name           string                 `json:"name"`
	Level          int                    `json:"level"`
	Count          int64                  `json:"count"`
	Type           domain.AchievementType `json:"type" enums:"hand,win"`
	LevelCount     int64                  `json:"levelCount"`
	NextLevelCount int64                  `json:"nextLevelCount"`
	NextLevelPrize int64                  `json:"nextLevelPrize"`
}

func ToAchievements(a *domain.AchievementStore) []*Achievement {
	result := HandCounterToAchievements(a.HandsCounter)
	for _, counter := range a.GetCounters() {
		result = append(result, CounterToAch(counter))
	}
	return result
}

func CounterToAch(c domain.Counter) *Achievement {
	if c == nil {
		return nil
	}
	return &Achievement{
		Name:           c.GetName(),
		Level:          c.GetLevel(),
		Count:          c.GetCount(),
		Type:           c.GetType(),
		LevelCount:     c.LevelCount(),
		NextLevelCount: c.NextLevelCount(),
		NextLevelPrize: c.NextLevelPrize(),
	}
}

func HandCounterToAchievements(hc *domain.HandsCounter) []*Achievement {
	counters := hc.OrderedCounters()
	var achs = make([]*Achievement, 0, len(counters))
	for _, counter := range counters {
		achs = append(achs, &Achievement{
			Name:           string(counter.Hand),
			Type:           domain.HandType,
			Level:          counter.Level,
			Count:          counter.Count,
			LevelCount:     counter.LevelCount(),
			NextLevelCount: counter.NextLevelCount(),
			NextLevelPrize: counter.NextLevelPrize(),
		})
	}
	return achs
}
