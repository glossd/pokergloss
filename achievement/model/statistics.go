package model

import "github.com/glossd/pokergloss/achievement/domain"

type Statistics struct {
	UserID       string `json:"userId"`
	GameCount    int64  `json:"gameCount"`
	WinPercent   int64  `json:"winPercent"`
	AllInPercent int64  `json:"allInPercent"`
	FoldPercent  int64  `json:"foldPercent"`
}

func ToStatistics(s *domain.Statistics) *Statistics {
	if s.GameCount == 0 {
		return &Statistics{
			UserID: s.UserID,
		}
	}
	return &Statistics{
		UserID:       s.UserID,
		GameCount:    s.GameCount,
		WinPercent:   s.WinCount * 100 / s.GameCount,
		AllInPercent: s.AllInCount * 100 / s.GameCount,
		FoldPercent:  s.FoldCount * 100 / s.GameCount,
	}
}
