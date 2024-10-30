package model

import "github.com/glossd/pokergloss/table/domain"

type Winner struct {
	Position int    `json:"position"`
	Chips    int64  `json:"chips"`
	HandRank string `json:"handRank"`
	Username string `json:"username"`
}

func ToWinners(t *domain.Table) *[]*Winner {
	var winners []*Winner
	for _, winner := range t.Winners {
		p, err := t.GetPlayer(winner.Position)
		var username string
		if err == nil {
			username = p.Username
		}
		winners = append(winners, ToWinner(winner, username))
	}
	return &winners
}

func ToWinner(w domain.Winner, username string) *Winner {
	return &Winner{
		Position: w.Position,
		Chips:    w.Chips,
		HandRank: w.HandRank,
		Username: username,
	}
}
