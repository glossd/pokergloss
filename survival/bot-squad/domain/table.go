package domain

import (
	"encoding/json"
	"github.com/pokerblow/poker"
	log "github.com/sirupsen/logrus"
	"sort"
)

type Table struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Size             int    `json:"size"`
	BigBlind         int64  `json:"bigBlind"`
	Status           string `json:"status"`
	DecidingPosition int    `json:"decidingPosition"`
	// Ordered, A-2.
	CommCards   Cards `json:"communityCards"`
	TotalPot    int64 `json:"totalPot"`
	MaxRoundBet int64 `json:"maxRoundBet"`

	// merge yourself
	Seats []*Seat `json:"seats"`

	// used to be taken from global var.
	// Being set once on init of the table
	UserPosition int
}

func NewTable(data []byte) (*Table, error) {
	var t Table
	err := json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

type Seat struct {
	Position int     `json:"position"`
	Blind    Blind   `json:"blind"`
	Player   *Player `json:"player"`
}

type Blind string

const (
	BigBlind   Blind = "bigBlind"
	SmallBlind Blind = "smallBlind"
	Dealer     Blind = "dealer"
	// if two players on a table, the dealer gets smallBlind
	DealerSmallBlind Blind = "dealerSmallBlind"
)

func (b Blind) IsDealer() bool {
	return b == Dealer || b == DealerSmallBlind
}

func (b Blind) IsSmallBlind() bool {
	return b == SmallBlind || b == DealerSmallBlind
}

func (b Blind) IsBigBlind() bool {
	return b == BigBlind
}

func (t *Table) SortCommunityCards() {
	sort.Slice(t.CommCards, func(i, j int) bool {
		return t.CommCards[i].FaceRank() > t.CommCards[j].FaceRank()
	})
}

func (t *Table) RankHoleCards(positions []int) {
	for _, pos := range positions {
		seat := t.Seats[pos]
		if seat.Player != nil {
			seat.Player.HoleCards.SortByFace()
			seat.Player.CardsRank = seat.Player.HoleCards.GetPatternRank()
		} else {
			log.Warnf("RankHoleCards no player at seat %d", pos)
		}
	}
}

func (t *Table) ResetPlayers() {
	for _, seat := range t.Seats {
		if seat.Player != nil {
			seat.Player.LastGameActionType = ""
			seat.Player.TotalRoundBet = 0
			seat.Player.HoleCards = nil
		}
	}
}

func (t *Table) GetPlayer(position int) *Player {
	return t.Seats[position].Player
}

func (t *Table) DecidingPlayer() *Player {
	return t.Seats[t.DecidingPosition].Player
}

func (t *Table) UserPlayer() *Player {
	return t.Seats[t.UserPosition].Player
}

type StackSize string

const (
	MediumStack StackSize = "medium"
	ShortStack  StackSize = "short"
	TinyStack   StackSize = "tiny"
)

func (t *Table) DecidingStackSize() StackSize {
	stack := t.DecidingPlayer().Stack
	switch {
	case stack > 35*t.BigBlind:
		return MediumStack
	case stack > 15*t.BigBlind:
		return ShortStack
	default:
		return TinyStack
	}
}

func (t *Table) IsDecidingBeforeUserPostFlop() bool {
	if t.UserPlayer().IsDealer() {
		return true
	}

	// Table Size > 2

	dealer := t.DealerPosition()
	deciding := t.DecidingPosition
	user := t.UserPosition
	var sortedPositions []int
	for i := dealer + 1; i < t.Size; i++ {
		sortedPositions = append(sortedPositions, i)
	}
	for i := 0; i <= dealer; i++ {
		sortedPositions = append(sortedPositions, i)
	}

	var decidingIdx int
	var userIdx int
	for i, position := range sortedPositions {
		if position == deciding {
			decidingIdx = i
		}
		if position == user {
			userIdx = i
		}
	}
	return decidingIdx < userIdx
}

func (t *Table) DealerPosition() int {
	for _, seat := range t.Seats {
		if seat.Blind == Dealer || seat.Blind == DealerSmallBlind {
			return seat.Position
		}
	}
	log.Errorf("Failed to find dealer position")
	return 0
}

func (t *Table) IsPreFlop() bool {
	return len(t.CommCards) == 0
}

func (t *Table) IsPostFlop() bool {
	return len(t.CommCards) >= 3
}

func (t *Table) IsFlop() bool {
	return len(t.CommCards) == 3
}

func (t *Table) IsTurn() bool {
	return len(t.CommCards) == 4
}

func (t *Table) IsRiver() bool {
	return len(t.CommCards) == 5
}

func (t *Table) MinRaiseChips() int64 {
	return 2*t.MaxRoundBet - t.DecidingPlayer().TotalRoundBet
}

func (t *Table) DecidingHand() poker.Hand {
	var cards []string
	for _, card := range t.CommCards {
		cards = append(cards, card.String())
	}
	for _, card := range t.DecidingPlayer().HoleCards {
		cards = append(cards, card.String())
	}
	if len(cards) < 5 {
		log.Errorf("DecidingHand has less than 5 cards, got=%d", len(cards))
		return poker.HighCard
	}
	r := poker.Eval(cards)
	return r.Hand
}

func (t *Table) IsTheUserFolded() bool {
	p := t.Seats[t.UserPosition].Player
	if p != nil {
		if p.LastGameActionType == FoldType {
			return true
		}
	}
	return false
}

func (t *Table) IsOthersMadeAllIn() bool {
	for _, seat := range t.Seats {
		if seat.Player != nil && seat.Position != t.DecidingPosition && seat.Player.LastGameActionType != FoldType {
			if seat.Player.LastGameActionType != AllInType {
				return false
			}
		}
	}
	return true
}

type EvalResult struct {
	Hand poker.Hand
	Rank int32
	Cards
}

func (t *Table) DecidingEval() EvalResult {
	var cards []string
	for _, card := range t.CommCards {
		cards = append(cards, card.String())
	}
	for _, card := range t.DecidingPlayer().HoleCards {
		cards = append(cards, card.String())
	}
	r := poker.Eval(cards)
	return EvalResult{
		Hand:  r.Hand,
		Rank:  r.Rank,
		Cards: toCards(r.Cards),
	}
}

// 0 - highest,
func (t *Table) flushWinRank(highHoleCard Card) int {
	var winRank = 0
	suitedComCards := t.communityCardsBySuit(highHoleCard.Suit())
	for _, face := range facesDesc {
		if face == highHoleCard.Face() {
			break
		}
		if suitedComCards.Contains(CardFrom(face, highHoleCard.Suit())) {
			continue
		}
		winRank++
	}
	return winRank
}

// 1 - best, 0 - worst
func (t *Table) flushWinRankage(highHoleCard Card) float64 {
	var winRank = 0
	suitedComCards := t.communityCardsBySuit(highHoleCard.Suit())
	for _, face := range facesDesc {
		if face == highHoleCard.Face() {
			break
		}
		if suitedComCards.Contains(CardFrom(face, highHoleCard.Suit())) {
			continue
		}
		winRank++
	}
	freeCards := 13 - len(suitedComCards)
	return float64(freeCards-winRank) / float64(freeCards)
}

func (t *Table) pocketPairRankage(pairFace Face) float64 {
	uCards := t.CommCards.UniqueByFace()
	for i, card := range uCards {
		if pairFace.Rank() > card.FaceRank() {
			return float64(len(uCards)+1-i) / float64(len(uCards)+1)
		}
	}
	// the worst
	return 0.2
}

func (t *Table) twoPairRankage(hc HoleCards) float64 {
	if hc.IsPair() {
		return t.pocketPairRankage(hc[0].Face())
	}

	if is, faces := t.CommCards.ContainsTwoPairs(); is {
		match := t.CommCards.FindFirstFaceMatch(hc)
		if match == UnknownFace {
			return 0.1 * t.highCardRankage(hc[0].Face())
		} else {
			if match.Rank() > faces[0].Rank() {
				return 1
			} else if match.Rank() > faces[1].Rank() {
				return 2.0 / 3
			} else {
				return 0.1 * t.highCardRankage(hc[0].Face())
			}
		}
	}

	if is, _ := t.CommCards.ContainsPair(); is {
		match := t.CommCards.FindFirstFaceMatch(hc)
		return t.pairRankageNoCheck(match)
	}

	combs := C(len(t.CommCards), 2)
	for i, card := range t.CommCards {
		if hc.ContainsFace(card.Face()) {
			for j, card := range t.CommCards[i+1:] {
				if card.Face() == hc[1].Face() {
					var combSkip = j
					if i > 0 {
						for k := i; k > 0; k-- {
							combSkip += len(t.CommCards) - k
						}
					}
					return float64(combs-combSkip) / float64(combs)
				}
			}
		}
	}
	return 0
}

// todo kicker
func (t *Table) pairRankage(hc HoleCards) float64 {
	if hc.IsPair() {
		return t.pocketPairRankage(hc[0].Face())
	}
	if is, _ := t.CommCards.ContainsPair(); is {
		return 0.1 * t.highCardRankage(hc[0].Face())
	}
	return t.pairRankageNoCheck(t.CommCards.FindFirstFaceMatch(hc))
}

func (t *Table) pairRankageNoCheck(pairFace Face) float64 {
	uCards := t.CommCards.UniqueByFace()
	for i, card := range uCards {
		if pairFace == card.Face() {
			return float64(len(uCards)-i) / float64(len(uCards))
		}
	}
	return 0
}

func (t *Table) setRankage(hc HoleCards) float64 {
	if is, _ := t.CommCards.ContainsSet(); is {
		return 0.1 * t.highCardRankage(hc[0].Face())
	}
	if hc.IsPair() {
		return 0.5 + 0.5*t.pocketPairRankage(hc[0].Face())
	} else {
		match := t.CommCards.FindFirstFaceMatch(hc)
		return t.pairRankageNoCheck(match)
	}
}

func (t *Table) highCardsRankage(hc HoleCards) float64 {
	// todo second card
	return t.highCardRankage(hc[0].Face())
}

func (t *Table) highCardRankage(face Face) float64 {
	uCards := t.CommCards.UniqueByFace()
	return float64(face.Rank()+1-len(uCards.AllLessThan(face))) / float64(13-len(uCards))
}

func (t *Table) straightRankage(hc HoleCards) float64 {
	cards := t.FindCardsOfStraight()
	if len(cards) < 3 {
		return 0
	}
	allGaps, areGapsStraight := cards.Gaps()
	if len(cards) == 3 {
		if allGaps == 2 {
			return 1
		}
		if allGaps == 1 {
			if cards[0].Face() == Ace {
				return 1
			}
			if hc.ContainsFace(cards[0].Face().inc()) {
				return 1
			} else {
				return 0.5
			}
		}
		if allGaps == 0 {
			if cards[0].Face() == Ace {
				return 1
			}
			if cards[0].Face() == King {
				if hc.ContainsFace(Ace) {
					return 1
				} else {
					return 0.5
				}
			}
			if hc.ContainsFace(cards[0].Face().inc().inc()) {
				return 1
			}
			if hc.ContainsFace(cards[0].Face().inc()) {
				return 0.9
			} else {
				return 0.45
			}
		}
	}

	if len(cards) == 4 {
		switch allGaps {
		case 4:
			// [K, T,9, 6]
			if hc.ContainsFace(cards[0].Face().dec()) {
				return 1
			} else {
				return 0.9
			}
		case 3:
			if hc.ContainsFace(cards[0].Face().decStraight()) {
				return 1
			} else {
				return 0.5
			}
		case 2:
			if cards[1].FaceRank()-3 == cards[2].FaceRank() {
				return 1
			}
			if areGapsStraight {
				if cards[0].FaceRank()-3 == cards[1].FaceRank() {
					if hc.ContainsFace(cards[0].Face().decStraight()) {
						return 1
					} else {
						return 0.5
					}
				} else {
					if hc.ContainsFace(cards[3].Face().decStraight()) {
						return 0.5
					} else {
						return 1
					}
				}
			} else {
				if cards[0].Face().decStraight() == cards[1].Face() {
					if cards[0].Face() == Ace {
						return 1
					}
					if hc.ContainsFace(cards[0].Face().inc()) {
						return 1
					} else {
						return 0.5
					}
				} else {
					if hc.ContainsFace(cards[0].Face().dec()) {
						return 1
					} else {
						return 0.5
					}
				}
			}
		case 1:
			if cards[0].Face() == Ace {
				return 1
			}
			if hc.ContainsFace(cards[0].Face().inc()) {
				return 1
			} else {
				return 0.5
			}
		case 0:
			if cards[0].Face() == Ace {
				return 1
			}
			if cards[0].Face() == King {
				if hc.ContainsFace(Ace) {
					return 1
				} else {
					return 0.5
				}
			}
			if hc.ContainsFace(cards[0].Face().inc().inc()) {
				return 1
			} else if hc.ContainsFace(cards[0].Face().inc()) {
				return 0.9
			} else {
				return 0.45
			}
		}
	}

	if len(cards) == 5 {
		switch allGaps {
		case 4:
			// [J, 9, 7, 5, 3]
			if hc.ContainsFace(cards[1].Face().inc()) {
				return 1
			} else if hc.ContainsFace(cards[2].Face().inc()) {
				return 2.0 / 3
			} else {
				return 1.0 / 3
			}
		case 3:
			if gaps, _ := cards[:3].Gaps(); gaps == 1 {
				if hc.ContainsFace(cards[0].Face().inc()) {
					return 1
				} else {
					// round
					return 0.5
				}
			} else {
				if hc.ContainsFace(cards[0].Face().decStraight()) {
					return 1
				} else {
					return 0.5
				}
			}
		case 2:
			if cards[0].FaceRank()-3 == cards[1].FaceRank() {
				if hc.ContainsFace(cards[0].Face().decStraight()) {
					return 1
				} else if hc.ContainsFace(cards[1].Face().inc()) {
					return 2.0 / 3
				} else {
					return 1.0 / 3
				}
			} else if gaps, _ := cards[:3].Gaps(); gaps == 2 {
				if hc.ContainsFace(cards[0].Face().dec()) {
					return 1
				} else if hc.ContainsFace(cards[1].Face().dec()) {
					return 2.0 / 3
				} else {
					return 1.0 / 3
				}
			} else if gaps == 1 {
				if cards[0].Face() == Ace {
					if cards[1].FaceRank()-1 == cards[2].FaceRank() {
						if hc.ContainsFace(cards[1].Face().inc()) {
							return 1
						} else {
							return 0.5
						}
					}
				}
				if hc.ContainsFace(cards[0].Face().inc()) {
					return 1
				} else {
					return 0.5
				}
			} else {
				if cards[0].Face() == Ace {
					return 1
				}
				if cards[1].Face() == King {
					if hc.ContainsFace(Ace) {
						return 1
					} else {
						return 0.5
					}
				}
				if hc.ContainsFace(cards[0].Face().inc().inc()) {
					// change of that is way lower
					return 1
				} else if hc.ContainsFace(cards[0].Face().inc()) {
					return 0.9
				} else {
					return 0.45
				}
			}
		case 1:
			if cards[0].Face() == Ace {
				return 1
			}
			if cards[0].Face() == King {
				if hc.ContainsFace(Ace) {
					return 1
				} else {
					return 0.5
				}
			}
			gapsOfFour, _ := cards[:4].Gaps()
			if hc.ContainsFace(cards[0].Face().inc().inc()) {
				return 1
			} else if hc.ContainsFace(cards[0].Face().inc()) {
				if gapsOfFour == 0 {
					return 0.8
				} else {
					return 0.95
				}
			} else {
				if gapsOfFour == 0 {
					return 0.45
				} else {
					return 0.6
				}
			}
		default:
			// case 0
			if cards[0].Face() == Ace {
				return 1
			}
			if cards[0].Face() == King {
				if hc.ContainsFace(Ace) {
					return 1
				} else {
					return 0.5
				}
			}
			if hc.ContainsFace(cards[0].Face().inc().inc()) {
				return 1
			}
			if hc.ContainsFace(cards[0].Face().inc()) {
				return 0.9
			} else {
				return 0.1
			}
		}
	}
	return 0
}

func (t *Table) FindCardsOfStraight() Cards {
	cards, idx := t.findThreeCardsStraight(0)
	if len(cards) < 3 {
		return nil
	}
	if idx == 0 {
		c, idx := t.findThreeCardsStraight(1)
		if len(c) >= 3 {
			if idx == 1 {
				c, _ := t.findThreeCardsStraight(2)
				if len(c) >= 3 {
					return t.CommCards
				} else {
					return t.CommCards[:4]
				}
			} else if idx == 2 {
				return t.CommCards
			}
		}
	} else if idx == 1 {
		c, _ := t.findThreeCardsStraight(2)
		if len(c) >= 3 {
			return t.CommCards[1:]
		}
	}
	return cards
}

func (t *Table) findThreeCardsStraight(fromIdx int) (Cards, int) {
	if len(t.CommCards)-fromIdx < 3 {
		return nil, 0
	}
	var beStraightCards = Cards{t.CommCards[fromIdx]}
	startStraightCard := t.CommCards[fromIdx]
	startStraightIdx := fromIdx
	straightGaps := 0
	var prevCard = startStraightCard
	for i, card := range t.CommCards[fromIdx+1:] {
		if prevCard.FaceRank() == card.FaceRank() {
			continue
		}
		if gaps := prevCard.Face().decStraight().Rank() - card.FaceRank(); gaps > 2 {
			startStraightCard = card
			startStraightIdx = i + 1
			prevCard = card
		} else {
			if straightGaps+gaps <= 2 {
				beStraightCards = append(beStraightCards, card)
				straightGaps += gaps
				prevCard = card
				if len(beStraightCards) >= 3 {
					break
				}
			} else {
				beStraightCards = nil
				straightGaps = 0
				startStraightCard = card
				prevCard = card
				startStraightIdx = i + 1
			}
		}
	}
	return beStraightCards, startStraightIdx
}

func (t *Table) communityCardsBySuit(suit Suit) (cards Cards) {
	for _, card := range t.CommCards {
		if card.Suit() == suit {
			cards = append(cards, card)
		}
	}
	return
}

func (t *Table) IsFlushDraw() bool {
	count, _ := t.CommCards.MaxSuitCount()
	return count >= 4
}

func (t *Table) IsStraightDraw() (is bool, containsGap bool) {
	return t.CommCards.IsStraightDraw()
}

// only for river
func (t *Table) IsStraight() bool {
	if len(t.CommCards) != 5 {
		return false
	}
	first := t.CommCards[0]
	if first.Face() == Ace {
		if t.CommCards[1:].IsGoingStraight(Five) {
			return true
		}
	}
	prevFaceRank := first.FaceRank()
	for _, card := range t.CommCards[1:] {
		if card.FaceRank() != prevFaceRank-1 {
			return false
		}
	}
	return true
}

func (t *Table) Pot() int64 {
	return t.TotalPot
}

func (t *Table) PotOdds() float64 {
	return float64(t.MaxRoundBet-t.DecidingPlayer().TotalRoundBet) / float64(t.Pot())
}

func (t *Table) betTimes() float64 {
	p := t.DecidingPlayer()
	maxBet := t.MaxRoundBet
	if t.MaxRoundBet > p.Stack {
		maxBet = p.Stack
	}
	return float64(maxBet) / float64(p.TotalRoundBet)
}

func (t *Table) HalfPot() int64 {
	return t.TotalPot / 2
}

func (t *Table) QuarterPot() int64 {
	return t.TotalPot / 2
}

func (t *Table) JSON() []byte {
	r, err := json.Marshal(t)
	if err != nil {
		log.Errorf("Failed to marshal table: %s", err)
		return nil
	}
	return r
}

func C(n, k int) int {
	return F(n) / (F(k) * F(n-k))
}

func F(n int) int {
	if n < 0 {
		return 1
	}
	res := 1
	for i := 1; i <= n; i++ {
		res *= i
	}
	return res
}
