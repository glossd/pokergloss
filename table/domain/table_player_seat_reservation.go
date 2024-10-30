package domain

import (
	log "github.com/sirupsen/logrus"
"github.com/glossd/pokergloss/auth/authid"
)

var ErrAlreadySitting = E("you already seating at this table")
var ErrCantCancel = E("you can't cancel, you seat is not in reservation")
var ErrSeatTaken = E("the seat is taken")

func (t *Table) ReserveSeat(position int, iden authid.Identity) error {
	switch t.Type {
	case SitngoType:
		return ErrNotAvailableInSitngo
	case MultiType:
		return ErrNotAvailableInMultiTable
	}
	seat, err := t.validateAddPlayer(position, iden)
	if err != nil {
		return err
	}

	seat.addPlayer(iden)

	return nil
}

func (t *Table) validateAddPlayer(position int, iden authid.Identity) (*Seat, error) {
	err := t.validatePosition(position)
	if err != nil {
		log.Errorf("User tried to seat at the table %s out of bound position %d, identity=%v", t.ID, position, iden)
		return nil, err
	}

	if t.ContainsPlayer(iden) {
		log.Warnf("User tried to seat again at the table %s, identity=%v", t.ID.Hex(), iden)
		return nil, ErrAlreadySitting
	}

	seat := t.GetSeatUnsafe(position)

	if seat.IsTaken() {
		log.Warnf("User tried to seat at the taken seat, tableID=%s, position=%d, identity=%s", t.ID.Hex(), position, iden)
		return nil, ErrSeatTaken
	}
	return seat, nil
}

func (t *Table) CancelSeatReservation(position int, iden authid.Identity) error {
	seat, err := t.GetSeatIdentified(position, iden)
	if err != nil {
		return err
	}
	p := seat.GetPlayer()
	if p.Status != PlayerReservedSeat {
		return ErrCantCancel
	}

	seat.RemovePlayer()
	return nil
}

func (t *Table) CancelSeatReservationTimeout(position int) error {
	seat, err := t.GetSeat(position)
	if err != nil {
		return err
	}

	if seat.IsFree() {
		log.Tracef("Cancel reservation of empty seat, tableID=%s", t.ID.Hex())
		return nil
	}
	p := seat.GetPlayer()
	if p.Status != PlayerReservedSeat {
		return ErrCantCancel
	}

	seat.RemovePlayer()
	return nil
}
