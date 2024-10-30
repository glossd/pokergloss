package model

import "github.com/glossd/pokergloss/achievement/domain"

type Exp struct {
	UserId           string `json:"userId"`
	Points           int64  `json:"points"`
	Level            int    `json:"level"`
	StartLevelPoints int64  `json:"startLevelPoints"`
	NextLevelPoints  int64  `json:"nextLevelPoints"`
}

func ToExp(e *domain.ExP) *Exp {
	return &Exp{
		UserId:           e.UserID,
		Points:           e.Points,
		Level:            e.Level,
		StartLevelPoints: e.StartLevelPoints(),
		NextLevelPoints:  e.NextLevelPoints(),
	}
}
