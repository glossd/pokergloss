package domain

import "github.com/pokerblow/poker"

func (b *Bot) flopAction(t *Table) Action {
	nuts := flopNuts(t.CommCards)
	switch nuts.Hand {
	case poker.StraightFlush:
		return b.flopStraightFlushNuts(t, nuts)
	case poker.FourOfAKind:
		return b.flopFourOfAKindNuts(t, nuts)
	case poker.Flush:
		return b.flopFlushNuts(t, nuts)
	case poker.Straight:
		return b.flopStraightNuts(t, nuts)
	case poker.ThreeOfAKind:
		return b.flopThreeOfAKindNuts(t, nuts)
	}
	return b.checkFold(t)
}

func (b *Bot) flopStraightFlushNuts(t *Table, nuts NutsResult) Action {
	p := t.DecidingPlayer()
	hand := t.DecidingHand()
	switch hand {
	case poker.StraightFlush:
		return b.doAction(t, 1)
	case poker.Flush:
		// both hole cards with suit of the flush
		rankage := t.flushWinRankage(p.HoleCards[0])
		return b.doAction(t, 0.4 +  0.5*rankage)
	case poker.Straight:
		if p.HoleCards.ContainsSuit(nuts.Suit) {
			return b.doAction(t, 0.9)
		}
		return b.doAction(t, 0.6*t.straightRankage(p.HoleCards))
	case poker.ThreeOfAKind:
		return b.doAction(t, 0.4 + 0.3*t.setRankage(p.HoleCards))
	case poker.TwoPair:
		return b.doAction(t, 0.3 + 0.3*t.twoPairRankage(p.HoleCards))
	case poker.Pair:
		return b.doAction(t, 0.2 + 0.3*t.pairRankage(p.HoleCards))
	case poker.HighCard:
		return b.doAction(t, 0.1*t.highCardsRankage(p.HoleCards))
	default:
		return b.checkFold(t)
	}
}

func (b *Bot) flopFourOfAKindNuts(t *Table, nuts NutsResult) Action {
	p := t.DecidingPlayer()
	hand := t.DecidingHand()
	if t.CommCards.AllSameFace() {
		if p.HoleCards.ContainsFace(t.CommCards[0].Face()) {
			return b.doAction(t, 1)
		} else if p.HoleCards.IsPair() {
			return b.doAction(t, 0.5 + 0.4*t.highCardRankage(p.HoleCards[0].Face()))
		} else {
			return b.doAction(t, 0.2*t.highCardsRankage(p.HoleCards))
		}
	}
	switch hand {
	case poker.FourOfAKind:
		// two cards in your hands and two on the board
		return b.doAction(t, 1)
	case poker.FullHouse:
		return b.doAction(t, 1)
	case poker.ThreeOfAKind:
		// two cards on the board, you got one
		return b.doAction(t, 0.7 + 0.1*t.highCardRankage(p.HoleCards.FindNotFace(nuts.Face).Face()))
	case poker.TwoPair:
		return b.doAction(t, 0.5*t.twoPairRankage(p.HoleCards))
	case poker.Pair, poker.HighCard:
		return b.doAction(t, 0.1*t.highCardRankage(p.HoleCards[0].Face()))
	}
	return b.checkFold(t)
}

func (b *Bot) flopFlushNuts(t *Table, nuts NutsResult) Action {
	p := t.DecidingPlayer()
	hand := t.DecidingHand()
	// board has three suited cards
	switch hand {
	case poker.Flush:
		return b.doAction(t, 0.7 + 0.3*t.flushWinRankage(p.HoleCards[0]))
	case poker.ThreeOfAKind:
		return b.doAction(t, 0.5 + 0.3*t.pairRankageNoCheck(p.HoleCards[0].Face()))
	case poker.TwoPair:
		return b.doAction(t, 0.3 + 0.4*t.twoPairRankage(p.HoleCards))
	case poker.Pair:
		return b.doAction(t, 0.2 + 0.4*t.pairRankage(p.HoleCards))
	case poker.HighCard:
		return b.doAction(t, 0.1*t.highCardRankage(p.HoleCards[0].Face()))
	default:
		return b.checkFold(t)
	}
}

func (b *Bot) flopStraightNuts(t *Table, nuts NutsResult) Action {
	p := t.DecidingPlayer()
	hand := t.DecidingHand()

	if p.HoleCards.AreSuited() {
		if t.CommCards.IsFlushDraw(p.HoleCards) {
			return b.doAction(t, 0.5)
		}
	}

	switch hand {
	case poker.Straight:
		return b.doAction(t, t.straightRankage(p.HoleCards))
	case poker.ThreeOfAKind:
		return b.doAction(t, 0.5 + 0.4*t.setRankage(p.HoleCards))
	case poker.TwoPair:
		return b.doAction(t, 0.4 + 0.4*t.twoPairRankage(p.HoleCards))
	case poker.Pair:
		return b.doAction(t, 0.2 + 0.4*t.pairRankage(p.HoleCards))
	case poker.HighCard:
		return b.doAction(t, 0.1*t.highCardRankage(p.HoleCards[0].Face()))
	default:
		return b.checkFold(t)
	}
}

func (b *Bot) flopThreeOfAKindNuts(t *Table, nuts NutsResult) Action {
	p := t.DecidingPlayer()
	hand := t.DecidingHand()

	switch hand {
	case poker.ThreeOfAKind:
		return b.doAction(t, t.setRankage(p.HoleCards))
	case poker.TwoPair:
		return b.doAction(t, 0.9 * t.twoPairRankage(p.HoleCards))
	case poker.Pair:
		return b.doAction(t, 0.6*t.pairRankage(p.HoleCards))
	case poker.HighCard:
		allCards := append(t.CommCards, p.HoleCards...)
		isStraightDraw, containsGap := allCards.IsStraightDraw()
		// board only has two suited
		if t.CommCards.IsFlushDraw(p.HoleCards) {
			if isStraightDraw {
				if containsGap {
					return b.doAction(t, 0.6)
				} else {
					return b.doAction(t, 0.7)
				}
			} else {
				return b.doAction(t, 0.5)
			}
		}
		if isStraightDraw {
			if containsGap {
				return b.doAction(t, 0.2)
			} else {
				return b.doAction(t, 0.4)
			}
		}
		return b.doAction(t, 0.1*t.highCardRankage(p.HoleCards[0].Face()))
	default:
		return b.checkFold(t)
	}
}

func (t *Table) flopPairFaceAndSingleFace() (pair Face, single Face) {
	if t.CommCards[0].Face() == t.CommCards[1].Face() {
		return t.CommCards[0].Face(), t.CommCards[2].Face()
	} else {
		return t.CommCards[1].Face(), t.CommCards[0].Face()
	}
}
