package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBotPreFlop_LoosenessOnPreFlopWithBigStack(t *testing.T) {
	action := bestBot().getAction(tableBotDealer(botData{cards: holeCards("Th", "3d"), stack: 129}))
	assert.EqualValues(t, Fold, action)
}

func TestBotPreFlop_BugRaisesWithTrashCardsOnPreflop(t *testing.T) {
	action := bestBot().getAction(tableBotBB(botData{cards: holeCards("6s", "9c"), stack: 50}))
	assert.EqualValues(t, Check, action)
}

func bestBot() *Bot {
	return &Bot{Position: 1, Looseness: 0.5, Aggression: 0.85}
}

type botData struct {
	cards HoleCards
	stack int64
}

func tableBotDealer(data botData) *Table {
	return &Table{
		MaxRoundBet: 2,
		BigBlind: 2,
		Seats: []*Seat{
			{Position: 0, Blind: BigBlind, Player: &Player{Position: 0, Blind: BigBlind, Stack: 98, TotalRoundBet: 2}},
			{Position: 1, Blind: DealerSmallBlind, Player: &Player{Position: 1, Blind: DealerSmallBlind, Stack: data.stack, TotalRoundBet: 1, HoleCards: data.cards, CardsRank: data.cards.GetPatternRank()}}},
		DecidingPosition: 1,
	}
}

func tableBotBB(data botData) *Table {
	return &Table{
		MaxRoundBet: 2,
		BigBlind: 2,
		Seats: []*Seat{
			{Position: 0, Blind: DealerSmallBlind, Player: &Player{Position: 0, Blind: DealerSmallBlind, Stack: 98, TotalRoundBet: 2}},
			{Position: 1, Blind: BigBlind, Player: &Player{Position: 1, Blind: BigBlind, Stack: data.stack, TotalRoundBet: 2, HoleCards: data.cards, CardsRank: data.cards.GetPatternRank()}}},
		DecidingPosition: 1,
	}
}

func holeCards(f string, s string) HoleCards {
	hc := HoleCards{Card(f), Card(s)}
	hc.SortByFace()
	return hc
}