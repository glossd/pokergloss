package domain

import (
	"github.com/pokerblow/poker"
	"math"
)

func (b *Bot) riverAction(t *Table) Action {
	nuts := turnOrRiverNuts(t.CommCards)
	switch nuts.Hand {
	case poker.StraightFlush:
		return b.riverStraightFlushNuts(t, nuts)
	case poker.FourOfAKind:
		return b.riverQuadsNuts(t, nuts)
	case poker.Flush:
		return b.riverFlushNuts(t, nuts)
	case poker.Straight:
		return b.riverStraightNuts(t, nuts)
	case poker.ThreeOfAKind:
		return b.riverThreeOfAKindNuts(t, nuts)
	default:
		return b.checkFold(t)
	}
}

func (b *Bot) riverStraightFlushNuts(t *Table, nuts NutsResult) Action {
	p := t.DecidingPlayer()
	eval := t.DecidingEval()

	isTableStraight := t.IsStraight()
	isTableFlush := t.CommCards.AllSameSuit()
	if isTableStraight && isTableFlush {
		highCard := t.CommCards[0]
		if highCard.Face() == Ace {
			return b.doAction(t, 1)
		}
		if p.HoleCards.Contains(CardFrom(highCard.Face().inc(), highCard.Suit())) {
			return b.doAction(t, 1)
		} else {
			return b.checkFold(t)
		}
	}


	switch eval.Hand {
	case poker.StraightFlush:
		if isTableFlush {
			if is, containsGap := t.CommCards.IsStraightDraw(); is {
				if containsGap {
					return b.doAction(t, 1)
				} else {
					return b.doAction(t, 0.95)
				}
			} else {
				// all two hole cards in the StraightFlush
				return b.doAction(t, 1)
			}
		} else {
			if is, containsGap := t.CommCards.FilterBySuit(nuts.Suit).IsStraightDraw(); is {
				if containsGap {
					return b.doAction(t, 1)
				} else {
					return b.doAction(t, 0.95)
				}
			} else {
				return b.doAction(t, 1)
			}
		}
	case poker.FourOfAKind:
		return b.doAction(t, 0.9)
	case poker.FullHouse:
		if is, face := t.CommCards.ContainsSet(); is {
			theCards := t.CommCards.AllExceptFace(face)
			if p.HoleCards.ContainsFace(theCards[0].Face()) {
				return b.doAction(t, 0.7)
			} else {
				return b.doAction(t, 0.5)
			}
		} else {
			return b.doAction(t, 0.9)
		}
	case poker.Flush:
		suitedCards := t.CommCards.FilterBySuit(nuts.Suit)
		if len(suitedCards) == 5 {
			rankage := t.flushWinRankage(p.HoleCards.FindFirstWithSuit(nuts.Suit))
			if rankage == 1 {
				return b.doAction(t, 0.8)
			} else if rankage > 0.7 {
				return b.doAction(t, 0.4 * rankage)
			} else {
				return b.checkFold(t)
			}
		}
		if len(suitedCards) == 4 {
			return b.doAction(t, 0.7*t.flushWinRankage(p.HoleCards.FindFirstWithSuit(nuts.Suit)))
		}
		// len(suitedCards) == 3
		return b.doAction(t, 0.7 + 0.25*t.flushWinRankage(p.HoleCards.FindFirstWithSuit(nuts.Suit)))
	case poker.Straight:

		var minusConfidence float64
		if t.IsFlushDraw() {
			minusConfidence = 0.3
		}
		return b.doAction(t, 0.4 + 0.3*t.straightRankage(p.HoleCards) - minusConfidence)
	case poker.ThreeOfAKind:
		return b.doAction(t, 0.4 + 0.2*t.setRankage(p.HoleCards) - b.riverTakeAway(t))
	case poker.TwoPair:
		return b.doAction(t, 0.3 + 0.3*t.twoPairRankage(p.HoleCards) - b.riverTakeAway(t))
	case poker.Pair:
		return b.doAction(t, 0.1 + 0.3*t.pairRankage(p.HoleCards) - 0.5*b.riverTakeAway(t))
	case poker.HighCard:
		return b.doAction(t, 0.1*t.highCardRankage(p.HoleCards[0].Face()))
	default:
		return b.checkFold(t)
	}
}

func (b *Bot) riverQuadsNuts(t *Table, nuts NutsResult) Action {
	p := t.DecidingPlayer()
	eval := t.DecidingEval()

	if is, face := t.CommCards.ContainsQuads(); is {
		if face == Ace {
			if t.CommCards.ContainsFace(King) {
				return b.doAction(t, 1)
			}
		} else {
			if t.CommCards.ContainsFace(Ace) {
				return b.doAction(t, 1)
			}
		}
		rankage := t.highCardRankage(p.HoleCards[0].Face())
		return b.doAction(t, math.Pow(rankage, rankage))
	}

	isTableFlushDraw := t.IsFlushDraw()

	switch eval.Hand {
	case poker.FourOfAKind:
		return b.doAction(t, 1)
	case poker.FullHouse:
		if is, face := t.CommCards.ContainsSet(); is {
			otherCards := t.CommCards.AllExceptFace(face)
			if p.HoleCards.IsPair() {
				if p.HoleCards[0].FaceRank() > otherCards[0].FaceRank() {
					return b.doAction(t, 0.9)
				} else if p.HoleCards[0].FaceRank() > otherCards[1].FaceRank() {
					return b.doAction(t, 0.6)
				} else {
					return b.doAction(t, 0.5)
				}
			}
			if p.HoleCards.ContainsFace(otherCards[0].Face()) {
				return b.doAction(t, 0.8)
			} else {
				return b.doAction(t, 0.5)
			}
		} else {
			if is, faces := t.CommCards.ContainsTwoPairs(); is {
				if p.HoleCards.ContainsFace(faces[0]) {
					return b.doAction(t, 0.9)
				} else {
					return b.doAction(t, 0.8)
				}
			} else {
				_, face := t.CommCards.ContainsPair()
				if t.CommCards[0].Face() == face {
					if p.HoleCards.ContainsFace(t.CommCards[2].Face()) {
						return b.doAction(t, 0.99)
					} else if p.HoleCards.ContainsFace(t.CommCards[3].Face()) {
						return b.doAction(t, 0.97)
					} else {
						return b.doAction(t, 0.95)
					}
				} else if t.CommCards[1].Face() == face {
					if p.HoleCards.ContainsFace(t.CommCards[0].Face()) {
						return b.doAction(t, 0.98)
					} else {
						return b.doAction(t, 0.95)
					}
				} else if t.CommCards[2].Face() == face {
					if p.HoleCards.ContainsFace(t.CommCards[0].Face()) {
						return b.doAction(t, 0.97)
					} else {
						return b.doAction(t, 0.94)
					}
				} else {
					if p.HoleCards.ContainsFace(t.CommCards[0].Face()) {
						return b.doAction(t, 0.96)
					} else {
						return b.doAction(t, 0.93)
					}
				}
			}
		}
	case poker.Flush:
		// board flush is impossible, since nuts is Quads
		if isTableFlushDraw {
			return b.doAction(t, 0.8 * t.flushWinRankage(p.HoleCards.FindFirstWithSuit(nuts.Suit)))
		} else {
			return b.doAction(t, 0.7 + 0.2*t.flushWinRankage(p.HoleCards.FindFirstWithSuit(nuts.Suit)))
		}
	case poker.Straight:
		var minusConfidence float64
		if isTableFlushDraw {
			minusConfidence = 0.3
		}
		// board straight is impossible, since nuts is Quads
		return b.doAction(t, 0.4 + 0.3*t.straightRankage(p.HoleCards) - minusConfidence)
	case poker.ThreeOfAKind:
		return b.doAction(t, 0.4 + 0.2*t.setRankage(p.HoleCards) - b.riverTakeAway(t))
	case poker.TwoPair:
		return b.doAction(t, 0.3 + 0.25*t.twoPairRankage(p.HoleCards) - b.riverTakeAway(t))
	case poker.Pair, poker.HighCard:
		return b.doAction(t, 0.1*t.highCardRankage(p.HoleCards[0].Face()))
	default:
		return b.checkFold(t)
	}
}

func (b *Bot) riverFlushNuts(t *Table, nuts NutsResult) Action {
	p := t.DecidingPlayer()
	eval := t.DecidingEval()
	if t.CommCards.AllSameSuit() {
		if p.HoleCards.ContainsSuit(nuts.Suit) {
			return b.doAction(t, t.highCardRankage(p.HoleCards.FindFirstWithSuit(nuts.Suit).Face()))
		} else {
			return b.checkFold(t)
		}
	}
	isTableFlushDraw := t.IsFlushDraw()

	switch eval.Hand {
	case poker.Flush:
		if isTableFlushDraw {
			return b.doAction(t, t.flushWinRankage(p.HoleCards.FindFirstWithSuit(nuts.Suit)))
		} else {
			return b.doAction(t, 0.7 + 0.3*t.flushWinRankage(p.HoleCards.FindFirstWithSuit(nuts.Suit)))
		}
	case poker.Straight:
		if isTableFlushDraw {
			return b.doAction(t, 0.2 + 0.3*t.straightRankage(p.HoleCards))
		} else {
			return b.doAction(t, 0.3 + 0.4*t.straightRankage(p.HoleCards))
		}
	case poker.ThreeOfAKind:
		return b.doAction(t, 0.8*t.setRankage(p.HoleCards) - b.riverTakeAway(t))
	case poker.TwoPair:
		return b.doAction(t, 0.7*t.twoPairRankage(p.HoleCards) - b.riverTakeAway(t))
	case poker.Pair:
		return b.doAction(t, 0.5*t.pairRankage(p.HoleCards) - b.riverTakeAway(t))
	case poker.HighCard:
		return b.doAction(t, 0.1*t.highCardRankage(p.HoleCards[0].Face()))
	default:
		return b.checkFold(t)
	}
}

func (b *Bot) riverStraightNuts(t *Table, nuts NutsResult) Action {
	p := t.DecidingPlayer()
	eval := t.DecidingEval()
	switch eval.Hand {
	case poker.Straight:
		return b.doAction(t, t.straightRankage(p.HoleCards))
	case poker.ThreeOfAKind:
		return b.doAction(t, 0.9*t.setRankage(p.HoleCards) - b.riverTakeAway(t))
	case poker.TwoPair:
		return b.doAction(t, 0.8*t.twoPairRankage(p.HoleCards) - b.riverTakeAway(t))
	case poker.Pair:
		return b.doAction(t, 0.6*t.pairRankage(p.HoleCards) - b.riverTakeAway(t))
	case poker.HighCard:
		return b.doAction(t, 0.1*t.highCardRankage(p.HoleCards[0].Face()))
	default:
		return b.checkFold(t)
	}
}

func (b *Bot) riverThreeOfAKindNuts(t *Table, nuts NutsResult) Action {
	p := t.DecidingPlayer()
	eval := t.DecidingEval()
	switch eval.Hand {
	case poker.ThreeOfAKind:
		return b.doAction(t, t.setRankage(p.HoleCards))
	case poker.TwoPair:
		return b.doAction(t, 0.9*t.twoPairRankage(p.HoleCards))
	case poker.Pair:
		return b.doAction(t, 0.8*t.pairRankage(p.HoleCards))
	case poker.HighCard:
		return b.doAction(t, 0.1*t.highCardRankage(p.HoleCards[0].Face()))
	default:
		return b.checkFold(t)
	}
}
