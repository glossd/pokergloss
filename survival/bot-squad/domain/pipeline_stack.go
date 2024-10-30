package domain

func (b *Bot) stackPipeline(a Action) Action {
	t := b.t
	p := t.DecidingPlayer()
	if !b.IsWeak {
		if a.IsAggressive() && p.Stack <= 5*t.BigBlind {
			if a.Chips > int64(3*float64(t.BigBlind)) {
				return AllIn
			}
		}
	}
	if a == Fold {
		if t.TotalPot > p.Stack*10 {
			return AllIn
		}
	}
	return a
}
