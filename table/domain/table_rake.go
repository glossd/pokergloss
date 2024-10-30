package domain

import conf "github.com/glossd/pokergloss/goconf"

func (t *Table) GetRake() Rake {
	return t.rake
}

func (t *Table) buildRake() Rake {
	if !t.IsCashType() {
		return Rake{}
	}
	if t.Status != GameEndTable {
		return Rake{}
	}
	// no flop, no drop
	if t.RoundType() == PreFlopRound {
		return Rake{}
	}
	if t.RakePercent == 0.0 {
		return Rake{}
	}
	winners := t.buildWinners()
	var posChips = make(map[int]int64)
	var allRakeChips int64
	for _, winner := range winners {
		userRake := t.computeRake(winner.Chips)
		if userRake > conf.Props.Table.MaxRake {
			userRake = conf.Props.Table.MaxRake
		}
		posChips[winner.Position] = userRake
		allRakeChips += userRake
	}
	return Rake{PositionChips: posChips, Chips: allRakeChips}
}

func (t *Table) computeRake(chips int64) int64 {
	return computeRakeFrom(chips, t.RakePercent)
}
