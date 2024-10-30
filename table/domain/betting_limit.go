package domain

import "math"

type BettingLimit string
const (
	NL BettingLimit = "NL"
	PL BettingLimit = "PL"
	ML BettingLimit = "ML"
	FL BettingLimit = "FL"
)

func (t *Table) CalcMaxAllowedBet(p *Player) int64 {
	max := t.BettingLimitChips() - p.TotalRoundBet
	if p.Stack < max {
		return p.Stack
	}
	return max
}

// don't care about totalRoundBet, just telling the game limit
func (t *Table) BettingLimitChips() int64 {
	switch t.BettingLimit {
	case FL:
		mrb := t.MaxRoundBet()
		if t.FixedLimitBet == 0 {
			return mrb + t.BigBlind
		}
		return mrb + t.FixedLimitBet
	case PL:
		return t.potLimitMaxBet()
	case NL:
		return math.MaxInt64
	case ML:
		if t.IsPreFlop() {
			return t.potLimitMaxBet()
		} else {
			return math.MaxInt64
		}
	default:
		return math.MaxInt64
	}
}

func (t *Table) potLimitMaxBet() int64 {
	return t.MaxRoundBet() + t.TotalPot
}
