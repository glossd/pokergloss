package domain

import (
	log "github.com/sirupsen/logrus"
)

func (b *Bot) preFlopAction(t *Table) Action {
	p := t.DecidingPlayer()

	b.confidence = p.cardsConfidence()

	if !b.loosenessPre(t) {
		log.Debug("PreFlop: looseness fold")
	}

	log.Debugf("PreFlop: confidence %v", b.confidence)

	if b.loosenessPre(t) {
		return b.aggressionPre(t)
	} else {
		if t.MaxRoundBet > t.BigBlind && isUserPreflopAggressive() {
			agro := userPreflopAggro()
			if agro > 0.8 {
				agro = 0.8
			}
			boost := agro - 0.2
			b.confidence += boost
			log.Debugf("PreFlop: aggro retry boost=%v, new confidence=%v", boost, b.confidence)
			if b.loosenessPre(t) {
				return Call
			}
		}
		if p.TotalRoundBet >= t.BigBlind && p.Stack < 3*t.BigBlind {
			if b.Aggression > 0.5 {
				return AllIn
			} else {
				return Call
			}
		}
	}

	if t.MaxRoundBet == t.BigBlind && t.DecidingPlayer().TotalRoundBet == t.MaxRoundBet {
		return Check
	}
	return Fold
}

// PreFlop
func preFlopCheckCall(t *Table) Action {
	if t.DecidingPlayer().TotalRoundBet >= t.MaxRoundBet {
		return Check
	} else {
		return Call
	}
}

func (b *Bot) loosenessPre(t *Table) bool {
	p := t.DecidingPlayer()
	if p.TotalRoundBet == t.MaxRoundBet {
		return true
	}

	coef := 1.0
	if t.DecidingStackSize() == ShortStack {
		coef = 0.75
	}
	if t.DecidingStackSize() == TinyStack {
		coef = 0.5
	}

	if p.TotalRoundBet > 0 {
		if t.DecidingStackSize() == MediumStack {
			coef = 1.75
		}
		if t.DecidingStackSize() == ShortStack {
			coef = 1.25
		}
		if t.DecidingStackSize() == TinyStack {
			coef = 1.5
		}
	}

	if t.MaxRoundBet > t.BigBlind && p.TotalRoundBet > 0 {
		betTimes := t.betTimes()
		switch {
		case betTimes >= 10:
			coef *= 0.4
		case betTimes >= 5:
			coef *= 0.6
		case betTimes >= 3:
			coef *= 0.8
		}
	}

	return b.confidence*coef >= b.Tightness()
}

func (b *Bot) aggressionPre(t *Table) Action {
	if b.RoundMadeActionCount > 2 + int(3*b.Aggression) {
		return preFlopCheckCall(t)
	}

	p := t.DecidingPlayer()
	isRaise := p.cardsConfidence() >= b.preFlopAggressionConfidence(t)
	if isRaise {
		log.Debugf("Preflop Raise: b.preFlopAggressionConfidence %v", b.preFlopAggressionConfidence(t))
	}
	if isRaise {
		if b.IsWeak {
			if p.CardsRank < 10 {
				return Raise(t.MinRaiseChips() + 2*t.MaxRoundBet)
			}
		}
		if p.cardsConfidence() > b.preFlopAggressionConfidence(t)*1.5 {
			return Raise(t.MinRaiseChips() + t.MaxRoundBet)
		}
		return Raise(t.MinRaiseChips())
	} else {
		return preFlopCheckCall(t)
	}
}
