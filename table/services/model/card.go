package model

import (
	"fmt"
	"github.com/glossd/pokergloss/table/domain"
)

type Card string

var SecretHoleCards = []Card{"Xx", "Xx"}

func holeCardsToCards(cards *domain.HoleCards) *[]Card {
	if cards != nil {
		return ToCards(cards.Get())
	}
	return nil
}

func ToCardsNil(cards []domain.Card) *[]Card {
	if len(cards) == 0 {
		return nil
	}
	return ToCards(cards)
}

func ToCards(cards []domain.Card) *[]Card {
	newCards := make([]Card, len(cards))
	for i, c := range cards {
		newCards[i] = toCard(c)
	}
	return &newCards
}

func toCard(c domain.Card) Card {
	return Card(fmt.Sprintf("%c%c", c.F, c.S))
}
