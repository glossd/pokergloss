package domain

import (
	"github.com/pokerblow/poker"
	"math"
)

func (b *Bot) turnAction(t *Table) Action {
	nuts := turnOrRiverNuts(t.CommCards)
	switch nuts.Hand {
	case poker.StraightFlush:
		return b.turnStraightFlushNuts(t, nuts)
	case poker.FourOfAKind:
		return b.turnFourOfAKindNuts(t, nuts)
	case poker.Flush:
		return b.turnFlushNuts(t, nuts)
	case poker.Straight:
		return b.turnStraightNuts(t, nuts)
	case poker.ThreeOfAKind:
		return b.turnThreeOfAKindNuts(t, nuts)
	}
	return b.checkFold(t)
}

func (b *Bot) turnStraightFlushNuts(t *Table, nuts NutsResult) Action {
	p := t.DecidingPlayer()
	eval := t.DecidingEval()
	var minusConfidence float64
	if t.CommCards.AllSameSuit() {
		minusConfidence = 0.5
	}
	switch eval.Hand {
	case poker.StraightFlush:
		if t.CommCards.AllSameSuit() {
			if is, containsGap := t.CommCards.IsStraightDraw(); is {
				if containsGap {
					return b.doAction(t, 1)
				} else {
					if t.CommCards[0].Face() == Ace {
						return b.doAction(t, 1)
					} else {
						if p.HoleCards.Contains(t.CommCards[0].IncFace()) {
							return b.doAction(t, 1)
						} else {
							return b.doAction(t, 0.95)
						}
					}
				}
			} else {
				// all hole cards in action
				return b.doAction(t, 1)
			}
		} else {
			// only three
			return b.doAction(t, 1)
		}
	case poker.FourOfAKind:
		// hole cards pair, pair in comm cards
		return b.doAction(t, 0.95)
	case poker.FullHouse:
		return b.doAction(t, 0.95)
	case poker.Flush:
		rankage := t.flushWinRankage(p.HoleCards.FindFirstWithSuit(nuts.Suit))
		if t.CommCards.AllSameSuit() {
			return b.doAction(t, 0.9*rankage)
		} else {
			return b.doAction(t, 0.5 + 0.4*rankage)
		}
	case poker.Straight:
		return b.doAction(t, 0.4 + 0.3*t.straightRankage(p.HoleCards) - minusConfidence*0.7)
	case poker.ThreeOfAKind:
		return b.doAction(t, 0.3 + 0.3*t.setRankage(p.HoleCards))
	case poker.TwoPair:
		return b.doAction(t, 0.5 + 0.25*t.twoPairRankage(p.HoleCards))
	case poker.Pair:
		if is, _ := t.CommCards.ContainsPair(); is {
			a, isDraw := b.turnDraw(t)
			if isDraw {
				return a
			}
			return b.doAction(t, 0.1*t.highCardsRankage(p.HoleCards))
		} else {
			return b.doAction(t, 0.1 + 0.2 * t.pairRankage(p.HoleCards))
		}
	case poker.HighCard:
		a, isDraw := b.turnDraw(t)
		if isDraw {
			return a
		}
		return b.doAction(t, 0.1*t.highCardsRankage(p.HoleCards))
	default:
		return b.checkFold(t)
	}
}

func (b *Bot) turnFourOfAKindNuts(t *Table, nuts NutsResult) Action {
	p := t.DecidingPlayer()
	eval := t.DecidingEval()

	if t.CommCards.AllSameFace() {
		rankage := t.highCardRankage(p.HoleCards[0].Face())
		return b.doAction(t, math.Pow(rankage, rankage))
	}

	if ok, face := t.CommCards.ContainsSet(); ok {
		if p.HoleCards.ContainsFace(face) {
			return b.doAction(t, 1)
		}
		theCard := t.CommCards.AllExceptFace(face)[0]
		if p.HoleCards.IsPair() {
			pairFace := p.HoleCards[0].Face()
			if pairFace == theCard.Face() {
				if pairFace.Rank() > face.Rank() {
					return b.doAction(t, 1)
				} else {
					return b.doAction(t, 0.7)
				}
			}
			if pairFace.Rank() > theCard.FaceRank() {
				return b.doAction(t, 0.7)
			} else {
				return b.doAction(t, 0.2)
			}
		} else {
			if p.HoleCards.Contains(theCard) {
				return b.doAction(t, 0.5)
			}
		}
	}

	// CommCards contains a pair or two pairs
	switch eval.Hand {
	case poker.FourOfAKind:
		// and bot has a pair
		if ok, faces := t.CommCards.ContainsTwoPairs(); ok {
			pairFace := p.HoleCards[0].Face()
			if faces[0] == pairFace {
				return b.doAction(t, 1)
			} else {
				return b.doAction(t, 0.99)
			}
		} else {
			return b.doAction(t, 1)
		}
	case poker.FullHouse:
		if ok, faces := t.CommCards.ContainsTwoPairs(); ok {
			if p.HoleCards.ContainsFace(faces[0]) {
				return b.doAction(t, 0.98)
			} else {
				return b.doAction(t, 0.85)
			}
		} else {
			_, face := t.CommCards.ContainsPair()
			if t.CommCards[0].Face() == face {
				if p.HoleCards[1].Face() == t.CommCards[2].Face() {
					return b.doAction(t, 0.95)
				} else {
					return b.doAction(t, 0.9)
				}
			} else if t.CommCards[1].Face() == face {
				if p.HoleCards.ContainsFace(t.CommCards[0].Face()) {
					return b.doAction(t, 0.95)
				} else {
					return b.doAction(t, 0.88)
				}
			} else {
				if p.HoleCards[0].Face() == t.CommCards[0].Face() {
					return b.doAction(t, 0.93)
				} else {
					return b.doAction(t, 0.85)
				}
			}
		}
	case poker.Flush:
		return b.doAction(t, 0.5 + 0.3*t.flushWinRankage(p.HoleCards[0]))
	case poker.Straight:
		// flush draw is impossible
		return b.doAction(t, 0.45 + 0.3+t.straightRankage(p.HoleCards))
	case poker.ThreeOfAKind:
		// check if flush/straight backdoor draw possible
		// hole cards is NOT PAIR
		return b.doAction(t, 0.45 + 0.3*t.setRankage(p.HoleCards))
	case poker.TwoPair:
		return b.doAction(t, 0.4*t.twoPairRankage(p.HoleCards))
	case poker.Pair, poker.HighCard:
		action, isDraw := b.turnDraw(t)
		if isDraw {
			return action
		}
		return b.doAction(t, 0.1*t.highCardsRankage(p.HoleCards))
	default:
		return b.checkFold(t)
	}
}

func (b *Bot) turnFlushNuts(t *Table, nuts NutsResult) Action {
	p := t.DecidingPlayer()
	eval := t.DecidingEval()
	var minusConfience float64
	if t.CommCards.AllSameSuit() {
		minusConfience = 0.5
	}
	switch eval.Hand {
	case poker.Flush:
		highCard := p.HoleCards.FindFirstWithSuit(nuts.Suit)
		rankage := t.flushWinRankage(highCard)
		if t.CommCards.AllSameSuit() {
			return b.doAction(t, math.Pow(rankage, rankage))
		} else {
			return b.doAction(t, 0.8 + 0.2*rankage)
		}
	case poker.Straight:
		// comm cards CAN'T BE SAME SUIT otherwise nuts would be straight flush
		// comm cards CAN'T CONTAIN PAIR otherwise nuts would be quads
		return b.doAction(t, 0.8 * t.straightRankage(p.HoleCards))
	case poker.ThreeOfAKind:
		return b.doAction(t, 0.5 + 0.25*t.setRankage(p.HoleCards) - minusConfience*0.75)
	case poker.TwoPair:
		// both hole cards have a pair
		return b.doAction(t, 0.4 + 0.3*t.twoPairRankage(p.HoleCards) - minusConfience*0.7)
	case poker.Pair:
		return b.doAction(t, 0.2 + 0.3*t.pairRankage(p.HoleCards) - minusConfience*0.5)
	case poker.HighCard:
		if p.HoleCards.ContainsSuit(nuts.Suit) {
			return b.doAction(t, 0.1 + 0.2 * t.highCardRankage(p.HoleCards.FindFirstWithSuit(nuts.Suit).Face()))
		}
		return b.doAction(t, 0.1*t.highCardsRankage(p.HoleCards))
	default:
		return b.checkFold(t)
	}
}

func (b *Bot) turnStraightNuts(t *Table, nuts NutsResult) Action {
	p := t.DecidingPlayer()
	eval := t.DecidingEval()
	isTableStraightDraw, containsGap := t.CommCards.IsStraightDraw()
	minusConfidence := 0.0
	if isTableStraightDraw {
		if containsGap {
			minusConfidence = 0.2
		} else {
			minusConfidence = 0.4
		}
	}
	// comm cards CAN'T CONTAIN PAIR otherwise nuts would be quads
	switch eval.Hand {
	case poker.Straight:
		return b.doAction(t, t.straightRankage(p.HoleCards))
	case poker.ThreeOfAKind:
		// only three comm cards could make a straight
		return b.doAction(t, 0.5 + 0.4*t.setRankage(p.HoleCards) - minusConfidence*0.9)
	case poker.TwoPair:
		// only three comm cards could make a straight
		return b.doAction(t, 0.4 + 0.3*t.twoPairRankage(p.HoleCards) - minusConfidence*0.7)
	case poker.Pair:
		return b.doAction(t, 0.2 + 0.3*t.pairRankage(p.HoleCards) - minusConfidence*0.5)
	case poker.HighCard:
		action, isDraw := b.turnDraw(t)
		if isDraw {
			return action
		}
		return b.doAction(t, 0.1*t.highCardsRankage(p.HoleCards))
	default:
		return b.checkFold(t)
	}
}

func (b *Bot) turnThreeOfAKindNuts(t *Table, nuts NutsResult) Action {
	p := t.DecidingPlayer()
	eval := t.DecidingEval()

	switch eval.Hand {
	case poker.ThreeOfAKind:
		return b.doAction(t, 0.9 + 0.1*t.setRankage(p.HoleCards))
	case poker.TwoPair:
		return b.doAction(t, 0.5 + 0.4*t.twoPairRankage(p.HoleCards))
	case poker.Pair:
		return b.doAction(t, 0.7*t.pairRankage(p.HoleCards))
	case poker.HighCard:
		action, isDraw := b.turnDraw(t)
		if isDraw {
			return action
		}
		return b.doAction(t, 0.1*t.highCardsRankage(p.HoleCards))
	default:
		return b.checkFold(t)
	}
}

func (b *Bot) turnDraw(t *Table) (Action, bool) {
	p := t.DecidingPlayer()
	allCards := append(t.CommCards, p.HoleCards...)
	isStraightDraw, containsGap := allCards.IsStraightDraw()
	if t.CommCards.IsFlushDraw(p.HoleCards) {
		// all two hole cards are suited
		if isStraightDraw {
			return b.doAction(t, 0.6), true
		} else {
			return b.doAction(t, 0.35+0.5*t.highCardRankage(p.HoleCards[0].Face())), true
		}
	}

	if isStraightDraw {
		if containsGap {
			return b.doAction(t, 0.1), true
		} else {
			return b.doAction(t, 0.2), true
		}
	}
	return Action{}, false
}
