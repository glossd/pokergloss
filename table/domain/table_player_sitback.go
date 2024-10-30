package domain

import "github.com/glossd/pokergloss/auth/authid"

var ErrNoSitBackInSittingOut = E("you can't sit back, you are not sitting out")

func (t *Table) SitBack(position int, iden authid.Identity) (bool, error) {
	player, err := t.GetPlayerIdentified(position, iden)
	if err != nil {
		return false, err
	}

	if player.Status != PlayerSittingOut {
		return false, ErrNoSitBackInSittingOut
	}

	if player.Stack == 0 {
		return false, E("you can't sit back, you don't have chips on table")
	}

	switch t.Type {
	case CashType:
		player.Status = PlayerReady
		if t.IsWaiting() {
			if t.IsEnoughPlayersForGame() {
				err := t.startFirstGame()
				if err != nil {
					return false, err
				}
				return true, nil
			}
		}
	case SitngoType, MultiType:
		player.Status = PlayerReady
		if t.IsWaiting() {
			err := t.startFirstGame()
			return true, err
		}
	}
	return false, nil
}
