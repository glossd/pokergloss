package model

import (
	"github.com/glossd/pokergloss/table/domain"
)

type EntrySitAndGo struct {
	UserId   string `json:"userId,omitempty"`
	Username string `json:"username,omitempty"`
	Picture  string `json:"picture,omitempty"`
	Position int    `json:"position"`
}

func ToEntries(from []*domain.EntrySitAndGo) []*EntrySitAndGo {
	result := make([]*EntrySitAndGo, 0, len(from))
	for _, e := range from {
		result = append(result, &EntrySitAndGo{
			UserId: e.UserId, Username: e.Username, Picture: e.Picture, Position: e.Position})
	}
	return result
}

type LobbySitAndGo struct {
	ID                       string                    `json:"id"`
	LobbyTable               LobbyTable                `json:"lobbyTable"`
	PlacesPaid               int                       `json:"placesPaid"`
	BuyIn                    int64                     `json:"buyIn"`
	BettingLimit             domain.BettingLimit       `json:"bettingLimit"`
	LevelIncreaseTimeMinutes int                       `json:"levelIncreaseTimeMinutes"`
	StartingChips            int64                     `json:"startingChips"`
	PrizePool                int64                     `json:"prizePool"`
	Entries                  []*EntrySitAndGo          `json:"entries"`
	Status                   domain.LobbyStatus        `json:"status"`
	TableID                  *string                   `json:"tableId,omitempty"`
	Prizes                   *[]domain.TournamentPrize `json:"prizes"`
	StartAt                  int64                     `json:"startAt"`
}

type LobbyTable struct {
	Name               string `json:"name"`
	Size               int    `json:"size"`
	BigBlind           int64  `json:"bigBlind"`
	DecisionTimeoutSec int    `json:"decisionTimeoutSec"`
}

func ToSNGLobbies(list []*domain.LobbySitAndGo) []*LobbySitAndGo {
	results := make([]*LobbySitAndGo, 0, len(list))
	for _, l := range list {
		results = append(results, ToSitAndGoLobby(l))
	}
	return results
}

func ToSitAndGoLobby(l *domain.LobbySitAndGo) *LobbySitAndGo {
	var tableId *string
	if l.TableID != nil {
		id := l.TableID.Hex()
		tableId = &id
	}

	return &LobbySitAndGo{
		ID: l.ID.Hex(),
		LobbyTable: LobbyTable{
			Name:               l.Name,
			Size:               l.Size,
			BigBlind:           l.BigBlind,
			DecisionTimeoutSec: int(l.DecisionTimeout.Seconds()),
		},
		PlacesPaid:               l.PlacesPaid,
		BuyIn:                    l.BuyIn,
		LevelIncreaseTimeMinutes: int(l.LevelIncreaseTime.Minutes()),
		StartingChips:            l.StartingChips,
		BettingLimit:             l.BettingLimit,
		StartAt:                  l.StartAt,
		Entries:                  ToEntries(l.Entries),
		Status:                   l.Status,
		TableID:                  tableId,
		Prizes:                   &l.Prizes,
		PrizePool:                l.ExpectedPrizePool(),
	}
}
