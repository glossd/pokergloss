package events

import (
	"github.com/glossd/pokergloss/table/domain"
)

type UserEvents struct {
	UserID string
	Events []*TableEvent
}

func BuildUserHoleCards(table *domain.Table) []UserEvents {
	var userEvents []UserEvents
	for _, player := range table.PlayingPlayersByGameType() {
		ue := UserEvents{UserID: player.UserId, Events: []*TableEvent{BuildTableHoleCards(table, player.UserId)}}
		userEvents = append(userEvents, ue)
	}
	return userEvents
}
