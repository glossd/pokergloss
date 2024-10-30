package events

import (
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/model"
)

func SeatsHoleCards(table *domain.Table, userID string) []*model.Seat {
	var seats []*model.Seat
	for _, player := range table.PlayingPlayersByGameType() {
		if player.UserId == userID {
			seats = append(seats, model.PlayerToSeat(player, model.ToPlayerOpenCards, model.NillifyStack))
		} else {
			seats = append(seats, model.PlayerToSeat(player, model.ToPlayerSecretCards, model.NillifyStack))
		}
	}
	return seats
}
