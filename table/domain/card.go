package domain

import (
	"fmt"
	"strings"
)

// Suit represents the suit of the card (spade, heart, diamond, club)
type Suit rune
// Constants for Suit ♠♥♦♣
const (
	Club    Suit = 'c'
	Diamond Suit = 'd'
	Heart   Suit = 'h'
	Spade   Suit = 's'
)

// Face represents the face of the card (ace, two...queen, king)
type Face rune
// Constants for Face
const (
	Two   Face = '2'
	Three Face = '3'
	Four  Face = '4'
	Five  Face = '5'
	Six   Face = '6'
	Seven Face = '7'
	Eight Face = '8'
	Nine  Face = '9'
	Ten   Face = 'T'
	Jack  Face = 'J'
	Queen Face = 'Q'
	King  Face = 'K'
	Ace   Face = 'A'
)

// Keep a card immutable
type Card struct {
	F Face
	S Suit
	Str string
}

func NewCard(f Face, s Suit) Card {
	var builder strings.Builder
	builder.WriteRune(rune(f))
	builder.WriteRune(rune(s))
	return Card{S: s, F: f, Str: builder.String()}
}

// Careful using it!!! Better use NewCard
func NewCardStr(str string) Card {
	return NewCard(Face(str[0]), Suit(str[1]))
}

// Careful using it!!!
func CardsStr(strs ...string) []Card {
	var cards []Card
	for _, str := range strs {
		cards = append(cards, NewCardStr(str))
	}
	return cards
}

func (c *Card) Suit() Suit {
	return c.S
}

func (c *Card) Face() Face {
	return c.F
}

func (c *Card) String() string {
	return fmt.Sprintf("%c%c", c.F, c.S)
}
