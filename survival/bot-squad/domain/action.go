package domain

type ActionType string

// check,bet,fold,call,raise,allIn
const (
	CheckType ActionType = "check"
	BetType   ActionType = "bet"
	FoldType  ActionType = "fold"
	CallType  ActionType = "call"
	RaiseType ActionType = "raise"
	AllInType ActionType = "allIn"
)

type Action struct {
	Type ActionType
	Chips int64
}

var (
	Fold  = Action{Type: FoldType}
	Call  = Action{Type: CallType}
	Check = Action{Type: CheckType}
	AllIn = Action{Type: AllInType}
)

func Bet(chips int64) Action {
	return Action{Type: BetType, Chips: chips}
}

func Raise(chips int64) Action {
	return Action{Type: RaiseType, Chips: chips}
}

func (a Action) IsAggressive() bool {
	return a.Type.IsAggressive()
}

func (at ActionType) IsAggressive() bool {
	return at == BetType || at == RaiseType
}
