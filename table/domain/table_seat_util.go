package domain

import (
	log "github.com/sirupsen/logrus"
"github.com/glossd/pokergloss/auth/authid"
)

// Use it wisely :)
func (t *Table) GetSeatUnsafe(position int) *Seat {
	seat, err := t.GetSeat(position)
	if err != nil {
		log.Errorf("Table position out of bound, tableID=%s, position=%d", t.ID, position)
		return nil
	}

	return seat
}

func (t *Table) GetSeatIdentified(position int, iden authid.Identity) (*Seat, error) {
	err := t.validateIdentityPosition(iden, position)
	if err != nil {
		return nil, err
	}
	return t.Seats[position], nil
}

func (t *Table) GetSeat(position int) (*Seat, error) {
	err := t.validatePosition(position)
	if err != nil {
		return nil, err
	}
	return t.Seats[position], nil
}
