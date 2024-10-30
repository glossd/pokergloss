package domain

import "fmt"

type ActionType string

// check,bet,fold,call,raise,allIn
const (
	Check ActionType = "check"
	Bet   ActionType = "bet"
	Fold  ActionType = "fold"
	Call  ActionType = "call"
	Raise ActionType = "raise"
	AllIn ActionType = "allIn"
)

type Action struct {
	Type ActionType
	Chips int64
}

var (
	FoldAction = Action{Type: Fold, Chips: 0}
	CallAction = Action{Type: Call, Chips: 0}
	CheckAction = Action{Type: Check, Chips: 0}
	AllInAction = Action{Type: AllIn, Chips: 0}
)

func (a Action) String() string {
	return fmt.Sprintf("Action{type: %s, chips: %d}", a.Type, a.Chips)
}

func (a Action) IsAggressive() bool {
	return a.Type == Bet || a.Type == Raise || a.Type == AllIn
}

func (a Action) HasChips() bool {
	return a.IsHandBet() || a.IsAutoBet()
}

func (a Action) IsHandBet() bool {
	return a.Type == Raise || a.Type == Bet
}

func (a Action) IsAutoBet() bool {
	return a.Type == Call || a.Type == AllIn
}

func (a Action) IsChipFree() bool {
	return a.Type == Fold || a.Type == Check
}

func BetAction(chips int64) Action {
	return Action{Type: Bet, Chips: chips}
}

func RaiseAction(chips int64) Action {
	return Action{Type: Raise, Chips: chips}
}
