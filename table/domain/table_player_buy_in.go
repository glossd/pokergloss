package domain

import (
	log "github.com/sirupsen/logrus"
"github.com/glossd/pokergloss/auth/authid"
)

// Returns true if game started or false if game is not started or already started
func (t *Table) BuyIn(stack int64, position int, iden authid.Identity) (bool, error) {
	if t.IsTournament() {
		return false, ErrNotAvailableInTournament
	}
	p, err := t.GetPlayerIdentified(position, iden)
	if err != nil {
		return false, err
	}

	isAllowed := p.Status == PlayerReservedSeat || p.IsBroke()

	if !isAllowed {
		return false, E("buy-in is not allowed, you already have stack of chips")
	}

	err = t.validateInitStackBoundary(stack)
	if err != nil {
		return false, err
	}

	t.setBuyIn(p, stack)

	if t.isReadyForGame() {
		err := t.startFirstGame()
		if err != nil {
			return false, err
		}
		log.Infof("Started new game in table %s", t.ID.Hex())
		return true, nil
	}

	return false, nil
}

func (t *Table) setBuyIn(p *Player, stack int64) {
	p.setInitStack(stack)
}
