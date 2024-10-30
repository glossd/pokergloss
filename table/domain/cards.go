package domain

// Cards served to players in the beginning of a poker game
type HoleCards struct {
	First  Card
	Second Card
}

func (hc *HoleCards) Get() []Card {
	return []Card{hc.First, hc.Second}
}

func (hc *HoleCards) AsStringArray() []string {
	if hc == nil {
		return nil
	}
	return []string{hc.First.String(), hc.Second.String()}
}

type CommunityCards struct {
	Flop  *Flop
	Turn  *Card
	River *Card
	newCards []Card
	isNewCardsAreDealtOnStartShowDown bool
}

func (cc CommunityCards) AvailableCards() []Card {
	var cards []Card
	if cc.Flop != nil {
		cards = append(cards, cc.Flop.AllCardsValue()...)
	}
	if cc.Turn != nil {
		cards = append(cards, *cc.Turn)
	}
	if cc.River != nil {
		cards = append(cards, *cc.River)
	}
	return cards
}

func (cc *CommunityCards) Available() []*Card {
	cards := cc.Flop.AllCards()
	if cc.Turn != nil {
		cards = append(cards, cc.Turn)
	}
	if cc.River != nil {
		cards = append(cards, cc.River)
	}
	return cards
}

func (cc *CommunityCards) IsFull() bool {
	return len(cc.Available()) == 5
}

func (cc *CommunityCards) setFlop(f Card, s Card, t Card) {
	cc.Flop = &Flop{
		FlopFirst:  f,
		FlopSecond: s,
		FlopThird:  t,
	}
	cc.newCards = append(cc.newCards, f, s, t)
}

func (cc *CommunityCards) setTurn(c Card) {
	cc.Turn = &c
	cc.newCards = append(cc.newCards, c)
}

func (cc *CommunityCards) setRiver(c Card) {
	cc.River = &c
	cc.newCards = append(cc.newCards, c)
}

func (cc *CommunityCards) reset() {
	cc.Flop = nil
	cc.Turn = nil
	cc.River = nil
	cc.newCards = nil
	cc.isNewCardsAreDealtOnStartShowDown = false
}

type Flop struct {
	FlopFirst  Card
	FlopSecond Card
	FlopThird  Card
}

func (f *Flop) AllCardsValue() []Card {
	return []Card{f.FlopFirst, f.FlopSecond, f.FlopThird}
}

func (f *Flop) AllCards() []*Card {
	if f != nil {
		return []*Card{&f.FlopFirst, &f.FlopSecond, &f.FlopThird}
	}
	return []*Card{}
}

type RoundType string

const (
	PreFlopRound RoundType = "preFlop" // no community cards are dealt
	FlopRound    RoundType = "flop"    // three cards on a table
	TurnRound    RoundType = "turn"    // four cards on a table
	RiverRound   RoundType = "river"   // five cards on a table
)

func (cc CommunityCards) RoundType() RoundType {
	if cc.Flop == nil {
		return PreFlopRound
	}
	if cc.Turn == nil {
		return FlopRound
	}
	if cc.River == nil {
		return TurnRound
	}
	return RiverRound
}

func (cc *CommunityCards) GetNewCards() []Card {
	return cc.newCards
}

func (cc *CommunityCards) AsStringArray() []string {
	cards := cc.AvailableCards()
	result := make([]string, 0, len(cards))
	for _, card := range cards {
		result = append(result, card.String())
	}
	return result
}

func (cc CommunityCards) RoundCommunityCards() []Card {
	switch cc.RoundType() {
	case FlopRound:
		return []Card{cc.Flop.FlopFirst, cc.Flop.FlopSecond, cc.Flop.FlopThird}
	case TurnRound:
		return []Card{*cc.Turn}
	case RiverRound:
		return []Card{*cc.River}
	}

	return []Card{}
}
