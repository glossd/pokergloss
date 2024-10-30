package domain

type Player struct {
	UserID string `json:"userId"`
	Position int `json:"position"`
	Blind `json:"blind"`
	TotalRoundBet int64 `json:"totalRoundBet"`
	Stack int64 `json:"stack"`
	LastGameActionType ActionType `json:"lastGameAction"`
	// Ordered, A-2.
	HoleCards `json:"cards"`
	CardsRank int
}

func (p *Player) cardsConfidence() float64 {
	return (169 - float64(p.CardsRank))/169
}
