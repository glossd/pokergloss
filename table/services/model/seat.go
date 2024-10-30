package model

import "github.com/glossd/pokergloss/table/domain"

type Seat struct {
	Position int          `json:"position"`
	Blind    domain.Blind `json:"blind"`
	Player   *Player      `json:"player"`
}

func ToSeats(seats []*domain.Seat, mapper PlayerMapper) []*Seat {
	var newSeats []*Seat
	for _, seat := range seats {
		newSeats = append(newSeats, ToSeat(seat, mapper))
	}
	return newSeats
}

func ToSeat(s *domain.Seat, mapper PlayerMapper) *Seat {
	return &Seat{Position: s.Position, Blind: s.Blind, Player: mapper(s.Player)}
}

func PlayersToSeats(players []*domain.Player, mapper PlayerMapper, afterMappers ...PlayerAfterMapper) []*Seat {
	seats := make([]*Seat, 0, len(players))
	for _, player := range players {
		seats = append(seats, PlayerToSeat(player, mapper, afterMappers...))
	}
	return seats
}

func PlayerToSeat(p *domain.Player, mapper PlayerMapper, afterMappers ...PlayerAfterMapper) *Seat {
	newP := mapper(p)
	for _, afterMapper := range afterMappers {
		afterMapper(newP)
	}
	return &Seat{Position: p.Position, Blind: p.Blind, Player: newP}
}

func PlayerToSeatTimeout(p *domain.Player, timeoutAt int64, mapper PlayerMapper) *Seat {
	player := mapper(p)
	player.TimeoutAt = &timeoutAt
	return &Seat{Position: p.Position, Blind: p.Blind, Player: player}
}

func EmptySeat(position int) *Seat {
	return &Seat{Position: position}
}
