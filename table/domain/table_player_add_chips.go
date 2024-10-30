package domain

import "github.com/glossd/pokergloss/auth/authid"

var ErrUseBuyInInsteadOfAddChips = E("you can't add chips, when you have 0 chips, you should use buy in")

func (t *Table) AddChips(chips int64, position int, iden authid.Identity) error {
	if t.IsSitngoType() || t.IsMultiType() {
		return ErrNotAvailableInSitngo
	}
	p, err := t.GetPlayerIdentified(position, iden)
	if err != nil {
		return err
	}

	if p.IsBroke() {
		return ErrUseBuyInInsteadOfAddChips
	}

	if t.MaxBuyInStack() < p.Stack+chips {
		return E("you tried to add too much chips, you can add %d at most", t.MaxBuyInStack()-p.Stack)
	}

	if t.IsInPlay() {
		p.addChipsOnGameStart(chips)
	} else {
		p.addChipsToStack(chips)
	}

	return nil
}
