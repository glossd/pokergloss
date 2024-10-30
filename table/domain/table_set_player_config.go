package domain

import "github.com/glossd/pokergloss/auth/authid"

func (t *Table) SetAutoMuck(autoMuck bool, position int, iden authid.Identity) error {
	p, err := t.GetPlayerIdentified(position, iden)
	if err != nil {
		return err
	}

	p.setAutoMuck(autoMuck)
	return nil
}

func (t *Table) SetAutoTopUp(autoTopUp bool, position int, iden authid.Identity) error {
	if t.IsTournament() {
		return ErrNotAvailableInTournament
	}

	p, err := t.GetPlayerIdentified(position, iden)
	if err != nil {
		return err
	}

	p.setAutoTopUp(autoTopUp)
	return nil
}

func (t *Table) SetAutoReBuy(autoReBuy bool, position int, iden authid.Identity) error {
	if t.IsTournament() {
		return ErrNotAvailableInTournament
	}

	p, err := t.GetPlayerIdentified(position, iden)
	if err != nil {
		return err
	}

	p.setAutoReBuy(autoReBuy)
	return nil
}
