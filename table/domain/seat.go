package domain

import (
	"github.com/glossd/pokergloss/auth/authid"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	log "github.com/sirupsen/logrus"
)

type Seat struct {
	// Position range from 0 to (Table.Size - 1).
	Position int
	Blind    Blind
	// Player is null when seat is empty.
	Player *Player

	// Optimistic locking for seat reservation.
	Version          int64
	wasPlayerRemoved bool
}

type Blind string

const (
	BigBlind   Blind = "bigBlind"
	SmallBlind Blind = "smallBlind"
	Dealer     Blind = "dealer"
	// if two players on a table, the dealer gets smallBlind
	DealerSmallBlind Blind = "dealerSmallBlind"
)

func newSeat(i int) *Seat {
	return &Seat{Position: i, Version: 0}
}

func NewTournamentSeat(position int, iden authid.Identity, stack int64) *Seat {
	s := newSeat(position)
	s.addPlayer(iden)
	s.GetPlayer().setInitStack(stack)
	return s
}

func (s *Seat) IsFree() bool {
	return s.Player == nil
}

func (s *Seat) IsTaken() bool {
	return !s.IsFree()
}

func (s *Seat) setBlind(blind Blind) {
	s.Blind = blind
	if s.IsTaken() {
		s.GetPlayer().Blind = blind
	} else {
		log.Errorf("Couldn't set blind, seat is free")
	}
}

func (s *Seat) addPlayer(iden authid.Identity) {
	s.Player = NewPlayer(iden, s.Position)
}

func (s *Seat) multiSitPlayer(p *Player) {
	p.multiPreviousPosition = p.Position
	s.Blind = ""
	p.Position = s.Position
	p.gameReset()
	p.Status = PlayerReady
	s.Player = p
}

func (s *Seat) RemovePlayer() {
	s.Player = nil
	s.wasPlayerRemoved = true
}

func (s *Seat) WasPlayerRemoved() bool {
	return s.wasPlayerRemoved
}

func (s *Seat) GetPlayer() *Player {
	return s.Player
}

func (s *Seat) reset() {
	s.Blind = ""
	if s.IsTaken() {
		s.GetPlayer().gameReset()
	}
}

func (s *Seat) ReservationTimeoutAt() int64 {
	if s.IsFree() {
		return -1
	}
	return timeutil.NowAdd(conf.Props.Table.SeatReservationTimeout)
}
