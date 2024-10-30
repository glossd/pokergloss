package domain

func (b *Bot) positionPipeline(a Action) Action {
	t := b.t
	if t.MaxRoundBet == 0 && t.DecidingPlayer().IsDealer() && a == Check {
		b.confidence += 0.1
		return b.aggressionPipeline(a)
	}
	return a
}
