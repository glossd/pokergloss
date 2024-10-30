package domain

import (
	log "github.com/sirupsen/logrus"
	"math"
)

type Bot struct {
	Position int
	Looseness float64 // [0, 1]
	Aggression float64 // [0, 1]
	IsWeak bool

	prevAction Action
	prevPotQuarters int64
	RoundMadeActionCount int

	t *Table
	confidence float64
}

func (b *Bot) Fear() float64 {
	return 1.0-b.Aggression
}

func (b *Bot) Tightness() float64 {
	return 1.0 - b.Looseness
}

func (b *Bot) GetConfidence() float64 {
	return b.confidence
}

func (b *Bot) GameReset() {
	b.prevAction = Action{}
	b.prevPotQuarters = 0
	b.RoundMadeActionCount = 0
}

func (b *Bot) preFlopAggressionConfidence(t *Table) float64 {
	rank := b.Aggression * 169 / 2
	p := t.DecidingPlayer()
	if t.MaxRoundBet > t.BigBlind {
		rank *= float64(p.TotalRoundBet)/float64(t.MaxRoundBet)
	}
	min := math.Min(rank, 169)
	passConfidence := (169-min)/169
	if (t.MaxRoundBet-p.TotalRoundBet)*5 > p.Stack && t.MaxRoundBet > t.BigBlind {
		passConfidence *= (t.betTimes()+1)/2
	}
	passConfidence = math.Min(passConfidence, 1)
	return passConfidence
}

func (b *Bot) GetAction(t *Table) Action {
	if t.IsTheUserFolded() {
		return Fold
	}

	action := b.getAction(t)

	action = b.allInAndMinChipsNormalization(t, action)

	if action.IsAggressive() && t.IsPostFlop() {
		b.prevPotQuarters = action.Chips / t.QuarterPot()
	}

	if action.IsAggressive() && t.IsOthersMadeAllIn() {
		action = Call
	}

	b.prevAction = action
	b.RoundMadeActionCount++
	return action
}

func (b *Bot) allInAndMinChipsNormalization(t *Table, action Action) Action {
	if action.IsAggressive() {
		if action.Type == BetType {
			if action.Chips < t.BigBlind {
				action.Chips = t.BigBlind
			}
		}
		if action.Type == RaiseType {
			if action.Chips < t.MinRaiseChips() {
				action.Chips = t.MinRaiseChips()
			}
		}
		if action == Call {
			p := t.DecidingPlayer()
			if p.Stack <= t.MaxRoundBet-p.TotalRoundBet {
				action = AllIn
			}
		}
		if action.Chips >= t.DecidingPlayer().Stack {
			action = AllIn
		}
	}

	if action == Call {
		p := t.DecidingPlayer()
		if p.Stack <= t.MaxRoundBet-p.TotalRoundBet {
			action = AllIn
		}
	}

	return action
}

func (b *Bot) getAction(t *Table) Action {
	p := t.DecidingPlayer()
	if p.Position != b.Position {
		log.Errorf("GetAction of bot who's not deciding")
		return Fold
	}

	switch {
	case t.IsPreFlop():
		return b.preFlopAction(t)
	case t.IsFlop():
		return b.flopAction(t)
	case t.IsTurn():
		return b.turnAction(t)
	case t.IsRiver():
		return b.riverAction(t)
	}

	return Fold
}

func (b *Bot) riverTakeAway(t *Table) float64 {
	var minusConfidence float64
	if t.IsFlushDraw() {
		minusConfidence += 0.225
	}
	if is, containsGap := t.IsStraightDraw(); is {
		if containsGap {
			minusConfidence += 0.1
		} else {
			minusConfidence += 0.2
		}
	}
	return minusConfidence
}
