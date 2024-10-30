package domain

import (
	"github.com/pokerblow/poker/cardb"
	"math"
	"sort"
)

type Cards []Card

func (cards Cards) Contains(c Card) bool {
	for _, card := range cards {
		if c == card {
			return true
		}
	}
	return false
}

func (cards Cards) ContainsBoth(c1 Card, c2 Card) bool {
	found := 0
	for _, card := range cards {
		if c1 == card || c2 == card {
			found++
			if found == 2 {
				return true
			}
		}
	}
	return false
}

func (cards Cards) ContainsPair() (bool, Face) {
	var faceMap = make(map[Face]int)
	for _, card := range cards {
		if _, ok := faceMap[card.Face()]; ok {
			return true, card.Face()
		}
		faceMap[card.Face()]++
	}
	return false, UnknownFace
}

func (cards Cards) ContainsTwoPairs() (bool, []Face) {
	var faceMap = make(map[Face]int)
	for _, card := range cards {
		faceMap[card.Face()]++
	}
	var found int
	var faces []Face
	for face, count := range faceMap {
		if count == 2 {
			found++
			if len(faces) == 0 {
				faces = append(faces, face)
			} else {
				if face.Rank() < faces[0].Rank() {
					faces = append(faces, face)
				} else {
					faces = []Face{face, faces[0]}
				}
			}
			if found == 2 {
				return true, faces
			}
		}
	}
	return false, nil
}

func (cards Cards) Gaps() (int, bool) {
	prevCard := cards[0]
	straightCards := cards
	if prevCard.Face() == Ace {
		straightCards = append(straightCards, prevCard)
	}
	var gaps int
	var areGapsStraight bool
	for _, card := range straightCards[1:] {
		gap := prevCard.Face().decStraight().Rank() - card.FaceRank()
		if gap == 2 {
			areGapsStraight = true
		}
		gaps += gap
		prevCard = card
	}
	return gaps, areGapsStraight
}

func (cards Cards) ContainsSet() (bool, Face) {
	var faceMap = make(map[Face]int)
	for _, card := range cards {
		if count, ok := faceMap[card.Face()]; ok && count == 2 {
			return true, card.Face()
		}
		faceMap[card.Face()]++
	}
	return false, UnknownFace
}

func (cards Cards) ContainsQuads() (bool, Face) {
	var faceMap = make(map[Face]int)
	for _, card := range cards {
		if count, ok := faceMap[card.Face()]; ok && count == 3 {
			return true, card.Face()
		}
		faceMap[card.Face()]++
	}
	return false, UnknownFace
}

func (cards Cards) AllLessThan(f Face) Cards {
	for i, card := range cards {
		if f.Rank() > card.FaceRank() {
			return cards[i:]
		}
	}
	return nil
}

func (cards Cards) AllExceptFace(f Face) (allExcept Cards) {
	for _, card := range cards {
		if card.Face() == f {
			allExcept = append(allExcept, card)
		}
	}
	return
}

func (cards Cards) ContainsFace(f Face) bool {
	for _, card := range cards {
		if card.Face() == f {
			return true
		}
	}
	return false
}

func (cards Cards) SuitCount(s Suit) (count int) {
	for _, card := range cards {
		if card.Suit() == s {
			count++
		}
	}
	return
}

func (cards Cards) FindFirstFaceMatch(hc HoleCards) Face {
	for _, card := range cards {
		if card.Face() == hc[0].Face() {
			return card.Face()
		}
		if card.Face() == hc[1].Face() {
			return card.Face()
		}
	}
	return UnknownFace
}

func (cards Cards) AllHaveSuit(s Suit) bool {
	for _, card := range cards {
		if card.Suit() != s {
			return false
		}
	}
	return true
}

func (cards Cards) AllSameSuit() bool {
	s := cards[0].Suit()
	for _, card := range cards[1:] {
		if card.Suit() != s {
			return false
		}
	}
	return true
}

func (cards Cards) AllSameFace() bool {
	f := cards[0].Face()
	for _, card := range cards[1:] {
		if card.Face() != f {
			return false
		}
	}
	return true
}

func (cards Cards) MaxSuitCount() (int, Suit) {
	var suitMap = make(map[Suit]int)
	for _, card := range cards {
		suitMap[card.Suit()]++
	}
	var max int
	var suit Suit
	for s, count := range suitMap {
		if max < count {
			max = count
			suit = s
		}
	}

	return max, suit
}

func (cards Cards) FilterBySuit(s Suit) Cards {
	var filtered Cards
	for _, card := range cards {
		if card.Suit() == s {
			filtered = append(filtered, card)
		}
	}
	return filtered
}

func (cards Cards) MaxFaceCount() (int, []Face) {
	var faceMap = make(map[Face]int)
	for _, card := range cards {
		faceMap[card.Face()]++
	}
	var max int
	var faces []Face
	for f, count := range faceMap {
		if max == count {
			faces = append(faces, f)
		}
		if max < count {
			max = count
			faces = append(faces, f)
		}
	}

	sortFaces(faces)

	return max, faces
}

// todo bug: [9s, 8s, 5s, 4s, 3s] will return 2 gaps
func (cards Cards) IsStraightBackdoorDraw() (is bool, gaps int) {
	return cards.isStraightBackdoorDraw(true)
}

// Don't sort for comm cards, they're sorted already
func (cards Cards) isStraightBackdoorDraw(sort bool) (bool, int) {
	var uCards = cards
	if sort {
		uCards = cards.UniqueSortedAscByFace()
	}
	if len(uCards) < 3 {
		return false, 0
	}

	isStraightOutOfThree := func(uCards []Card) bool {
		lastFaceRank := uCards[2].FaceRank()
		if uCards[2].Face() == Ace {
			lastFaceRank = -1
		}
		return lastFaceRank > uCards[0].FaceRank()-5
	}

	for i := 0; i <= len(uCards)-3; i++ {
		set := uCards[i:i+3]
		if isStraightOutOfThree(set) {
			lastFaceRank := set[2].FaceRank()
			if set[2].Face() == Ace {
				lastFaceRank = -1
			}
			return true, set[0].FaceRank() - 2 - lastFaceRank
		}
	}

	if cards[0].Face() == Ace {
		return append(cards[1:], cards[0]).isStraightBackdoorDraw(false)
	}

	return false, 0
}

func (cards Cards) IsStraightDraw() (is bool, containsGap bool) {
	count, startFace, endFace := cards.GoingStraightCount()
	if count >= 4 {
		return true, false
	}
	if count == 3 {
		if dec, ok := startFace.dec().DecrementedStraight(); ok {
			if cards.ContainsFace(dec) {
				return true, true
			}
		}
		if inc, ok := endFace.inc().Incremented(); ok {
			if cards.ContainsFace(inc) {
				return true, true
			}
		}
	}
	if count == 2 {
		uniqueCards := cards.UniqueSortedDescCardsByFace()
		if len(uniqueCards) == 4 {
			if uniqueCards[0].FaceRank() == uniqueCards[1].FaceRank()+1 ||
				uniqueCards[2].FaceRank() == uniqueCards[3].FaceRank()+1 {
				return true, true
			}
		}
		if len(uniqueCards) == 5 {
			var endFaceIdx int
			for i, card := range uniqueCards {
				if card.Face() == endFace {
					endFaceIdx = i
				}
			}
			if len(uniqueCards) > endFaceIdx+2 {// means more than two
				count, _, _ := uniqueCards[endFaceIdx+1:].GoingStraightCount()
				if count > 2 {
					return true, true
				}
			}
		}

	}
	return false, false
}

func (cards Cards) IsFlushDraw(hc HoleCards) bool {
	if hc.AreSuited() {
		if cards.SuitCount(hc[0].Suit()) == 2 {
			return true
		}
	} else {
		if cards.SuitCount(hc[0].Suit()) == 3 {
			return true
		}
		if cards.SuitCount(hc[1].Suit()) == 3 {
			return true
		}
	}
	return false
}

func (cards Cards) IsGoingStraight(to Face) bool {
	for _, card := range cards {
		if card.Face() != to {
			return false
		}
		to.DecrementedStraight()
	}
	return true
}

// count minimum number is 1.
// if there are two pairs, it returns the highest pair
func (cards Cards) GoingStraightCount() (count int, startFace Face, endFace Face) {
	sorted := cards.UniqueSortedAscByFace()
	var maxNum int
	var gotNum int
	var gotStartFace Face
	var prevFaceRank = sorted[0].FaceRank()
	for _, card := range sorted[1:] {
		if card.FaceRank() == prevFaceRank+1 {
			if gotNum == 0 {
				gotStartFace = card.Face()
			}
			gotNum++
			if gotNum >= maxNum {
				maxNum = gotNum
				startFace = gotStartFace
				endFace = card.Face()
			}
		} else {
			gotNum = 0
		}
		prevFaceRank = card.FaceRank()
	}
	return maxNum+1, startFace, endFace
}

func (cards Cards) StraightGaps() (gaps int, highFace Face) {
	sorted := cards.UniqueSortedAscByFace()

	if len(sorted) < 3 {
		return -1, UnknownFace
	}

	var prevFaceRank = sorted[0].FaceRank()
	for _, card := range sorted[1:] {
		fr := card.FaceRank()
		gaps += fr - prevFaceRank + 1
		prevFaceRank = fr
	}

	return gaps, sorted[len(sorted)-1].Face()
}

func (cards Cards) UniqueByFace() (res Cards) {
	prevCard := cards[0]
	res = append(res, prevCard)
	for _, card := range cards[1:] {
		if prevCard.Face() != card.Face() {
			res = append(res, card)
		}
		prevCard = card
	}
	return
}

func (cards Cards) UniqueSortedAscByFace() Cards {
	return cards.UniqueSortedCardsByFace(false)
}

func (cards Cards) UniqueSortedDescCardsByFace() Cards {
	return cards.UniqueSortedCardsByFace(true)
}

func (cards Cards) UniqueSortedCardsByFace(ask bool) Cards {
	tmp := make([]Card, len(cards))
	copy(tmp, cards)
	sort.Slice(tmp, func(i, j int) bool {
		if ask {
			return tmp[i].FaceRank() < tmp[j].FaceRank()
		} else {
			return tmp[i].FaceRank() > tmp[j].FaceRank()
		}
	})
	var filtered Cards
	prevFace := tmp[0].Face()
	for _, card := range tmp[1:] {
		if prevFace != card.Face() {
			filtered = append(filtered, card)
		}
		prevFace = card.Face()
	}
	return tmp
}

func (cards Cards) AtLeastNHaveSuit(n int, s Suit) bool {
	var found int
	for _, card := range cards {
		if card.Suit() == s {
			found++
			if found >= n {
				return true
			}

		}
	}
	return false
}

func toCards(c []cardb.Card) Cards {
	var result = make([]Card, 0, len(c))
	for _, card := range c {
		result = append(result, Card(card))
	}
	return result
}

// from A to 2
func sortFaces(faces []Face) {
	sort.Slice(faces, func(i, j int) bool {
		return faces[i].Rank() > faces[j].Rank()
	})
}

func maxInt(ints ...int) int {
	var m = math.MinInt32
	for _, n := range ints {
		if m < n {
			m = n
		}
	}
	return m
}
