package db

import (
	"github.com/glossd/pokergloss/table/domain"
	"go.mongodb.org/mongo-driver/bson"
)

func SeatUpdates(table *domain.Table) []bson.E {
	updates := make([]bson.E, 0, len(table.Seats))
	for _, seat := range table.Seats {
		if seat.IsTaken() || seat.WasPlayerRemoved() {
			updates = append(updates, SeatUpdate(seat)...)
		} else {
			updates = append(updates, seatBlind(seat))
		}
	}
	return updates
}

func SeatUpdate(seat *domain.Seat) []bson.E {
	var updates []bson.E
	updates = append(updates, seatPlayer(seat))
	updates = append(updates, seatBlind(seat))
	return updates
}

// for simple sit-back
func PlayerUpdateStatus(seat *domain.Seat) []bson.E {
	var updates []bson.E
	updates = append(updates, bson.E{Key: SeatDbPath(seat.Position) + ".player.status", Value: seat.Player.Status})
	return updates
}

func seatPlayer(seat *domain.Seat) bson.E {
	return bson.E{Key: SeatDbPath(seat.Position) + ".player", Value: seat.Player}
}

func seatBlind(seat *domain.Seat) bson.E {
	return bson.E{Key: SeatDbPath(seat.Position) + ".blind", Value: seat.Blind}
}

func SeatVersion(seat *domain.Seat) bson.E {
	return bson.E{Key: SeatDbPath(seat.Position) + ".version", Value: seat.Version}
}

func IncSeatVersion(seat *domain.Seat) bson.E {
	return bson.E{Key: SeatDbPath(seat.Position) + ".version", Value: seat.Version + 1}
}
