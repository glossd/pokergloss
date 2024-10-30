package domain

import "github.com/pokerblow/poker"

type NutsResult struct {
	poker.Hand
	// For flush.
	Suit
	// For straight.
	Face Face

	// second pair in comm cards
	Face2 Face

	Gaps int // Straight
}

func flopNuts(cardsSlice []Card) NutsResult {
	cards := Cards(cardsSlice)
	suitCount, theMaxSuit := cards.MaxSuitCount()
	faceCount, facesOfMaxCount := cards.MaxFaceCount()
	isStraight, _ := cards.IsStraightBackdoorDraw()

	isFlush := suitCount >= 3

	if isFlush && isStraight {
		return NutsResult{
			Hand: poker.StraightFlush,
			Suit: theMaxSuit,
		}
	}

	if faceCount > 1 {
		return NutsResult{Hand: poker.FourOfAKind, Face: facesOfMaxCount[0]}
	}

	if isFlush {
		return NutsResult{
			Hand: poker.Flush,
			Suit: theMaxSuit,
		}
	}

	if isStraight {
		return NutsResult{Hand: poker.Straight}
	}
	return NutsResult{Hand: poker.ThreeOfAKind, Face: facesOfMaxCount[0]}
}

func turnOrRiverNuts(commCards []Card) NutsResult {
	cards := Cards(commCards)
	suitCount, theMaxSuit := cards.MaxSuitCount()
	isFlush := suitCount >= 3
	if isFlush {
		suitCards := cards.FilterBySuit(theMaxSuit)
		if is, _ := suitCards.IsStraightBackdoorDraw(); is {
			return NutsResult{Hand: poker.StraightFlush, Suit: theMaxSuit}
		}
	}
	if is, face := cards.ContainsPair(); is {
		return NutsResult{Hand: poker.FourOfAKind, Face: face}
	}

	if isFlush {
		return NutsResult{Hand: poker.Flush, Suit: theMaxSuit}
	}

	isStraight, _ := cards.IsStraightBackdoorDraw()
	if isStraight {
		return NutsResult{Hand: poker.Straight}
	}
	return NutsResult{Hand: poker.ThreeOfAKind}
}

