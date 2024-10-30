package domain


type Hand string
const (
	HighCard Hand = "High Card"
	Pair Hand = "Pair"
	TwoPair Hand = "Two Pair"
	ThreeOfAKind Hand = "Three of a Kind"
	Straight Hand = "Straight"
	Flush Hand = "Flush"
	FullHouse Hand = "Full House"
	FourOfAKind Hand = "Four of a Kind"
	StraightFlush Hand = "Straight Flush"
)

var hands = map[Hand]struct{} {
	HighCard:{},
	Pair:{},
	TwoPair:{},
	ThreeOfAKind:{},
	Straight:{},
	Flush:{},
	FullHouse:{},
	FourOfAKind:{},
	StraightFlush:{},
}

func joinWithHands(t AssignmentType) (a []*AssignmentUnit) {
	for hand := range hands {
		a = append(a, NewWithHand(t, hand))
	}
	return
}

type TableRound string
const (
	PreFlop TableRound = "preFlop"
	Flop TableRound = "flop"
	Turn TableRound = "turn"
	River TableRound = "river"
)
func (r TableRound) GetName() string {
	switch r {
	case PreFlop: return "PreFlop"
	case Flop: return "Flop"
	case Turn: return "Turn"
	case River: return "River"
	default: return ""
	}
}
var rounds = []TableRound{"", PreFlop, Flop, Turn, River}

func joinWithRound(t AssignmentType, number int) (a []*AssignmentUnit) {
	for _, r := range rounds {
		a = append(a, NewWithNumberAndRound(t, number, r))
	}
	return
}

