package domain

import (
	log "github.com/sirupsen/logrus"
"github.com/glossd/pokergloss/auth/authid"
"math/rand"
)

var Algo algorithms = &RealAlgo{}

type algorithms interface {
	ChooseDealer([]*Seat) *Seat
	RandomAvailableOneCard(t *Table) Card
	RandomAvailableTwoCards(t *Table) (f Card, s Card)
	RandomAvailableThreeCards(t *Table) (f Card, s Card, th Card)
	ShuffleUsers(users []authid.Identity)
}

type RealAlgo struct{}

func (a *RealAlgo) ShuffleUsers(users []authid.Identity) {
	rand.Shuffle(len(users), func(i, j int) {
		users[i], users[j] = users[j], users[i]
	})
}

func (a *RealAlgo) RandomAvailableOneCard(t *Table) Card {
	available := ReturnAvailableCards(t.AllCards())
	return available[rand.Intn(len(available))]
}

func (a *RealAlgo) RandomAvailableTwoCards(t *Table) (f Card, s Card) {
	available := ReturnAvailableCards(t.AllCards())
	firstIdx := rand.Intn(len(available) - 1)
	f = available[firstIdx]
	available[firstIdx] = available[len(available)-1]
	available = available[:len(available)-1]

	secondIdx := rand.Intn(len(available) - 1)
	s = available[secondIdx]
	available[secondIdx] = available[len(available)-1]

	return
}

func (a *RealAlgo) RandomAvailableThreeCards(t *Table) (f Card, s Card, th Card) {
	available := ReturnAvailableCards(t.AllCards())

	firstIdx := rand.Intn(len(available) - 1)
	f = available[firstIdx]
	available[firstIdx] = available[len(available)-1]
	available = available[:len(available)-1]

	secondIdx := rand.Intn(len(available) - 1)
	s = available[secondIdx]
	available[secondIdx] = available[len(available)-1]
	available = available[:len(available)-1]

	thirdIdx := rand.Intn(len(available) - 1)
	th = available[thirdIdx]
	available[thirdIdx] = available[len(available)-1]

	return
}

// https://mindyourdecisions.com/blog/2010/10/26/choosing-the-dealer-in-poker-is-dealing-to-the-first-ace-a-fair-system/
// It's saying that dealer gets chosen randomly
func (a *RealAlgo) ChooseDealer(seats []*Seat) *Seat {
	dealerIdx := rand.Intn(len(seats))
	dealerSeat := seats[dealerIdx]
	return dealerSeat
}

type MockAlgo struct {
	queueOfFistCards []Card
	dealerPos        int
	skip             int
}

func (a *MockAlgo) ShuffleUsers(users []authid.Identity) {}

func NewMockAlgo(queueOfFistCards []Card) (*MockAlgo, error) {
	setOfCards := make(map[Card]struct{}, len(queueOfFistCards))
	for _, c := range queueOfFistCards {
		setOfCards[c] = struct{}{}
	}

	if len(setOfCards) != len(queueOfFistCards) {
		return nil, E("cards are not unique")
	}

	return &MockAlgo{queueOfFistCards: queueOfFistCards}, nil
}

func NewMockCards(cards ...string) *MockAlgo {
	algo, err := NewMockAlgo(CardsStr(cards...))
	if err != nil {
		log.Fatalf("Failed create mock: %s", err)
	}
	return algo
}

func NewMockCardsSkip(skip int, cards ...string) *MockAlgo {
	algo, err := NewMockAlgo(CardsStr(cards...))
	if err != nil {
		log.Fatalf("Failed create mock: %s", err)
	}
	algo.skip = skip
	return algo
}

func NewMockAlgoMultiGame(queuesOfCards ...[]Card) (*MockAlgo, error) {
	for _, queue := range queuesOfCards {
		setOfCards := make(map[Card]struct{}, len(queue))
		for _, c := range queue {
			setOfCards[c] = struct{}{}
		}

		if len(setOfCards) != len(queue) {
			return nil, E("cards are not unique")
		}
	}

	var allQueue []Card
	for _, cards := range queuesOfCards {
		allQueue = append(allQueue, cards...)
	}

	return &MockAlgo{queueOfFistCards: allQueue}, nil
}

func NewMockFull(dealerPos int, queuesOfCards ...[]Card) *MockAlgo {
	mock, err := NewMockAlgoMultiGame(queuesOfCards...)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	mock.dealerPos = dealerPos
	return mock
}

func (a *MockAlgo) randomAvailableCardV2(unavailable []Card) Card {
	if a.skip > 0 {
		a.skip--
		return ReturnAvailableCards(unavailable)[0]
	}
	if len(a.queueOfFistCards) > 0 {
		card := a.queueOfFistCards[0]
		a.queueOfFistCards = a.queueOfFistCards[1:]
		return card
	}
	return ReturnAvailableCards(unavailable)[0]
}

func (a *MockAlgo) RandomAvailableOneCard(t *Table) Card {
	return a.randomAvailableCardV2(t.AllCards())
}

func (a *MockAlgo) RandomAvailableTwoCards(t *Table) (f Card, s Card) {
	allCards := t.AllCards()
	f = a.randomAvailableCardV2(allCards)
	s = a.randomAvailableCardV2(append(allCards, f))
	return
}

func (a *MockAlgo) RandomAvailableThreeCards(t *Table) (f Card, s Card, th Card) {
	allCards := t.AllCards()
	f = a.randomAvailableCardV2(allCards)
	s = a.randomAvailableCardV2(append(allCards, f))
	th = a.randomAvailableCardV2(append(allCards, f, s))
	return
}

func (a *MockAlgo) ChooseDealer(seats []*Seat) *Seat {
	for _, s := range seats {
		if s.Position == a.dealerPos {
			return s
		}
	}
	return seats[0]
}
