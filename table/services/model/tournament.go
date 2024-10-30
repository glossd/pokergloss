package model

import (
	"github.com/glossd/pokergloss/table/domain"
)

type TournamentAttrs struct {
	LobbyID         string                   `json:"lobbyId"`
	Name            string                   `json:"name"`
	StartAt         int64                    `json:"startAt"`
	LevelIncreaseAt int64                    `json:"levelIncreaseAt"`
	NextSmallBlind  int64                    `json:"nextSmallBlind"`
	Prizes          []domain.TournamentPrize `json:"prizes"`
	MarketPrize     *MarketPrize             `json:"marketPrize,omitempty"`
}

func toTournamentAttrs(attrs domain.TournamentAttributes) *TournamentAttrs {
	if attrs.LevelIncreaseAt == 0 {
		return nil
	}
	return &TournamentAttrs{
		LevelIncreaseAt: attrs.LevelIncreaseAt,
		LobbyID:         attrs.LobbyID.Hex(),
		Name:            attrs.Name,
		StartAt:         attrs.StartAt,
		Prizes:          attrs.Prizes,
		NextSmallBlind:  attrs.NextSmallBlind,
		MarketPrize:     toMarketPrize(attrs.MarketPrize),
	}
}

type MultiAttrs struct{}
