package domain

import (
	log "github.com/sirupsen/logrus"
	"math/rand"
)

func (b *Bot) bluffPipeline(a Action) Action {
	t := b.t

	if b.IsWeak {
		return a
	}

	if a == Check && t.MaxRoundBet == 0 && t.DecidingPlayer().IsDealer() {
		if !anyAggressiveAction(getLastPostflopUserActions(len(t.CommCards)-2)) {
			if rand.Float64() > 0.9 {
				log.Infof("Making Dealer Bluff")
				return Bet(t.Pot())
			}
		}
	}

	if a.Type == BetType && t.MaxRoundBet == 0 &&
		b.confidence >= 0.9 &&
		isUserPostflopAggressive() && t.IsDecidingBeforeUserPostFlop() {
		if rand.Float64() > 0.8 {
			log.Infof("Making Check Raise, hc=%v, cc=%v", t.DecidingPlayer().HoleCards, t.CommCards)
			return Check
		}
	}

	return a
}
