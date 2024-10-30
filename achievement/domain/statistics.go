package domain

type Statistics struct {
	UserID string `bson:"_id"`
	GameCount int64
	WinCount int64
	AllInCount int64
	FoldCount int64
}

func NewStatistics(userID string) *Statistics {
	return &Statistics{UserID: userID}
}

func (ps *Statistics) Update(end *GameEnd) {
	player, ok := end.PlayersMap[ps.UserID]
	if !ok {
		return
	}
	ps.GameCount++
	switch player.LastAction {
	case "allIn":
		ps.AllInCount++
	case "fold":
		ps.FoldCount++
	}
	_, ok = end.WinnersMap[ps.UserID]
	if ok {
		ps.WinCount++
	}
}
