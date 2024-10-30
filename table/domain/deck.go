package domain

// Do not mutate the deck! In golang array isn't immutable by nature. httpS://stackoverflow.com/a/13137568/10160865
var Deck = [52]Card{
	NewCard(Two, Club), NewCard(Two, Diamond), NewCard(Two, Heart), NewCard(Two, Spade),
	NewCard(Three, Club), NewCard(Three, Diamond), NewCard(Three, Heart), NewCard(Three, Spade),
	NewCard(Four, Club), NewCard(Four, Diamond), NewCard(Four, Heart), NewCard(Four, Spade),
	NewCard(Five, Club), NewCard(Five, Diamond), NewCard(Five, Heart), NewCard(Five, Spade),
	NewCard(Six, Club), NewCard(Six, Diamond), NewCard(Six, Heart), NewCard(Six, Spade),
	NewCard(Seven, Club), NewCard(Seven, Diamond), NewCard(Seven, Heart), NewCard(Seven, Spade),
	NewCard(Eight, Club), NewCard(Eight, Diamond), NewCard(Eight, Heart), NewCard(Eight, Spade),
	NewCard(Nine, Club), NewCard(Nine, Diamond), NewCard(Nine, Heart), NewCard(Nine, Spade),
	NewCard(Ten, Club), NewCard(Ten, Diamond), NewCard(Ten, Heart), NewCard(Ten, Spade),
	NewCard(Jack, Club), NewCard(Jack, Diamond), NewCard(Jack, Heart), NewCard(Jack, Spade),
	NewCard(Queen, Club), NewCard(Queen, Diamond), NewCard(Queen, Heart), NewCard(Queen, Spade),
	NewCard(King, Club), NewCard(King, Diamond), NewCard(King, Heart), NewCard(King, Spade),
	NewCard(Ace, Club), NewCard(Ace, Diamond), NewCard(Ace, Heart), NewCard(Ace, Spade),
}

func (c Card) DeckIndex() int {
	var faceIdx int
	switch c.F {
	case Two:
		faceIdx = 0
	case Three:
		faceIdx = 4
	case Four:
		faceIdx = 8
	case Five:
		faceIdx = 12
	case Six:
		faceIdx = 16
	case Seven:
		faceIdx = 20
	case Eight:
		faceIdx = 24
	case Nine:
		faceIdx = 28
	case Ten:
		faceIdx = 32
	case Jack:
		faceIdx = 36
	case Queen:
		faceIdx = 40
	case King:
		faceIdx = 44
	case Ace:
		faceIdx = 48
	}

	idx := faceIdx
	switch c.S {
	case Diamond:
		idx += 1
	case Heart:
		idx += 2
	case Spade:
		idx += 3
	}

	return idx
}

func ReturnAvailableCards(unavailable []Card) []Card {
	unavailableIndexes := make(map[int]struct{})
	for _, card := range unavailable {
		unavailableIndexes[card.DeckIndex()] = struct{}{}
	}
	available := make([]Card, 52 - len(unavailableIndexes))
	availableI := 0
	for i, card := range Deck {
		if _, ok := unavailableIndexes[i]; !ok {
			available[availableI] = card
			availableI++
		}
	}
	return available
}


// used this to generate above deck code:
//func main() {
//	faces := []string{"Two","Three","Four","Five","Six","Seven","Eight","Nine","Ten","Jack","Queen","King","Ace"}
//	suits := []string{"Club", "Diamond", "Heart", "Spade"}
//
//	for _, F := range faces {
//		for _, s := range suits {
//			fmt.Printf("{F: %s, S: %s}, ", F, s)
//		}
//		fmt.Println()
//	}
//}
