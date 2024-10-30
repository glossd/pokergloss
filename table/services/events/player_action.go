package events

import (
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/model"
)

func BuildActionsOfPlayers(t *domain.Table) []*TableEvent {
	var evts []*TableEvent
	for i, p := range t.MadeActionPlayerPositions() {
		eventType := PlayerMadeAction
		if i == 0 && t.WasActionTimeout() {
			eventType = TimeToDecideTimeout
		}

		totalRoundBet := p.TotalRoundBet
		if t.IsNewRound() {
			// totalRoundBet 0 will be sent with event newBettingRound
			totalRoundBet = p.LastGameBet
		}
		stack := p.Stack
		if t.IsGameEnd() {
			stack = t.GetPlayerStackBeforeResult(p)
		}
		afterMapper := func(p *model.Player) {
			p.TotalRoundBet = &totalRoundBet
			p.Stack = &stack
		}

		maxBet := t.MaxRoundBet()
		blc := t.BettingLimitChips()
		evts = append(evts, &TableEvent{Type: eventType, Payload: M{
			"table": model.Table{
				MaxRoundBet:       &maxBet,
				BettingLimitChips: &blc,
				TotalPot:          &t.TotalPot,
				DecidingPosition:  &model.NegativeOne,
				Seats:             []*model.Seat{model.PlayerToSeat(p, model.ToPlayerNoCardsNotDeciding, afterMapper)},
			},
		}})
	}

	return evts
}
