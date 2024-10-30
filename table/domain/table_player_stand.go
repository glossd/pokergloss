package domain

import "github.com/glossd/pokergloss/auth/authid"

func (t *Table) Stand(position int, iden authid.Identity) (*Player, error) {
	player, err := t.GetPlayerIdentified(position, iden)
	if err != nil {
		return nil, err
	}

	err = t.stand(player)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (t *Table) StandTimeout(position int) error {
	player, err := t.GetPlayer(position)
	if err != nil {
		return err
	}
	return t.stand(player)
}

func (t *Table) stand(p *Player) error {
	switch t.Type {
	case MultiType:
		return ErrNotAvailableInMultiTable
	case SitngoType:
		return ErrNotAvailableInSitngo
	}
	if t.isStandFold(p) {
		p.isStandFolded = true
		err := t.doMakeAction(p, FoldAction)
		if err != nil {
			return err
		}
		t.nullifyPlayer(p)
		return nil
	}
	if t.isAllowedLeaveRightAway(p.Position) {
		t.removePlayer(p)
	} else {
		p.IsLeaving = true
	}
	return nil
}

func (t *Table) isStandFold(p *Player) bool {
	return t.IsRingType() && t.IsPlaying() && t.DecidingPosition == p.Position
}

func (t *Table) removePlayer(p *Player) {
	if t.IsSitngoType() || t.IsMultiType() {
		p.Stack = 0
	}
	t.GetSeatUnsafe(p.Position).RemovePlayer()
}

func (t *Table) isAllowedLeaveRightAway(position int) bool {
	p, err := t.GetPlayer(position)
	if err != nil {
		return false
	}

	switch t.Type {
	case SitngoType, MultiType:
		return false
	case CashType:
		switch p.Status {
		case PlayerSittingOut, PlayerReservedSeat:
			return true
		case PlayerReady:
			switch t.Status {
			case PlayingTable, ShowdownTable, WaitingTable:
				return true
			case GameEndTable:
				return false
			}
		case PlayerPlaying:
			if t.Status != GameEndTable && p.LastGameAction == Fold {
				return true
			}
			return false
		}
	}

	return false
}
