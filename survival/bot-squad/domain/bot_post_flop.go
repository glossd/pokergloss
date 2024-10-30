package domain

import (
	log "github.com/sirupsen/logrus"
	"math"
)

// PostFlop
func (b *Bot) checkFold(t *Table) Action {
	return b.doAction(t, 0)
}

func (b *Bot) doAction(t *Table, confidence float64) Action {
	if confidence > 1 {
		confidence = 1
	}
	if confidence < 0 {
		confidence = 0
	}

	b.confidence = confidence
	b.t = t

	log.Debugf("PostFlop: confidence=%v", b.confidence)

	return b.bluffPipeline(b.userTrackingPipeline(b.stackPipeline(b.positionPipeline(b.aggressionPipeline(b.loosenessPipeline())))))
}

func (b *Bot) loosenessPipeline() Action {
	t := b.t
	confidence := b.confidence
	loosenessConfidenceBoost := math.Min(0, b.Looseness-0.7)
	if t.MaxRoundBet == 0 {
		return Check
	} else {
		isSuperConfident := confidence >= (0.9 + 0.1*b.Tightness())
		isTinyBet := t.PotOdds() < 0.07
		if isSuperConfident || isTinyBet || (confidence + loosenessConfidenceBoost) >= t.PotOdds() {
			return Call
		} else {
			return Fold
		}
	}
}

func (b *Bot) aggressionPipeline(a Action) Action {
	if a == Fold {
		return a
	}
	t := b.t
	confidence := b.confidence
	if t.MaxRoundBet > 0 {
		//isReRaise := p.TotalRoundBet > 0
		if confidence == 1 && b.RoundMadeActionCount < 2 + int(3*b.Aggression) {
			return Raise(t.MinRaiseChips())
		}
		return Call
	} else {
		if confidence > 0.01 && confidence >= b.Fear() {
			if confidence < b.checkOverBetConfidence() {
				return Check
			}

			var quarters int64
			if b.IsWeak {
				if confidence < 0.4 {
					return Bet(t.BigBlind)
				}
				if confidence > 0.5+0.3*b.Fear() {
					quarters = int64(math.Sqrt(confidence) * 12)
				} else {
					quarters = int64(confidence * 4)
				}
			} else {
				if t.IsTurn() && t.IsRiver() {
					if b.prevPotQuarters > 0 {
						return Bet(t.QuarterPot() * b.prevPotQuarters)
					}
				}
				quarters = int64(confidence * 4)
			}
			return Bet(t.QuarterPot() * quarters)
		} else {
			return Check
		}
	}
}

func (b *Bot) checkOverBetConfidence() float64 {
	return 0.2 + 0.3*b.Fear()
}

