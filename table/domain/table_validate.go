package domain

import (
	log "github.com/sirupsen/logrus"
"github.com/glossd/pokergloss/auth/authid"
)


var (
	ErrPositionNotYours = E("table position is not yours")
	ErrBuyInNotEnough   = E("this amount of chips is not enough for buy-in")
	ErrBuyInTooMuch     = E("this amount of chips is too much for buy-in")
)

func (t *Table) validateInitStackBoundary(stack int64) error {
	if stack > t.MaxBuyInStack() {
		log.Warnf("User tried to put stack more than max allowed, stack=%d, max stack=%d", stack, t.MaxBuyInStack())
		return ErrBuyInTooMuch
	}
	if stack < t.MinBuyInStack() {
		log.Warnf("User tried to put stack less than min allowed, stack=%d, min stack=%d", stack, t.MinBuyInStack())
		return ErrBuyInNotEnough
	}
	return nil
}

func (t *Table) validatePosition(position int) error {
	if position < 0 || position >= len(t.Seats) {
		return E("no such position in the table %s, position=%d", t.ID, position)
	}
	return nil
}

func (t *Table) validateIdentityPosition(iden authid.Identity, position int) error {
	err := t.validatePosition(position)
	if err != nil {
		return err
	}
	player := t.GetPlayerUnsafe(position)
	if player == nil {
		log.Errorf("User tried something on empty seat, tableID=%s, position=%d identity=%s", t.ID, position, iden)
		return ErrPositionNotYours
	}
	if player.UserId != iden.UserId {
		log.Errorf("User tried something on position belonging to another user, tableID=%s, position=%d identity=%s", t.ID, position, iden)
		return ErrPositionNotYours
	}
	return nil
}
