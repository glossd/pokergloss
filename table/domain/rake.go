package domain

import conf "github.com/glossd/pokergloss/goconf"

type Rake struct {
	Chips         int64
	PositionChips map[int]int64
}

func (r Rake) Of(pos int) int64 {
	if r.PositionChips == nil {
		return 0
	}
	return r.PositionChips[pos]
}

func computeRakeFrom(chips int64, rakePercent float64) int64 {
	return int64(float64(chips) * rakePercent)
}

func GetRakePercent() float64 {
	return conf.Props.Table.RakePercent
}

func computeTournamentFeeFrom(chips int64, feePercent float64) int64 {
	return int64(float64(chips) * feePercent)
}

func GetTournamentFeePercent() float64 {
	return conf.Props.Tournament.FeePercent
}
