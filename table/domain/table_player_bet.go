package domain

func(t *Table) bet(p *Player, chips int64) error {
	err := p.betChips(chips)
	if err != nil {
	    return err
	}
	t.TotalPot += chips
	return nil
}

// If player don't have enough chips for blind, he should `bet` with what he got
func (t *Table) betBlind(p *Player) {
	switch p.Blind {
	case BigBlind:
		bet, _ := p.betChipsForBlind(t.BigBlind)
		t.TotalPot += bet
	case SmallBlind:fallthrough
	case DealerSmallBlind:
		bet, _ := p.betChipsForBlind(t.SmallBlind)
		t.TotalPot += bet
	}
}
