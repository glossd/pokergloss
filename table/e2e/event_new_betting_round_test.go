package e2e

import (
	"github.com/glossd/pokergloss/table/domain"
	"testing"
)

func TestFlop(t *testing.T) {
	t.Cleanup(cleanUp)
	domain.Algo = &domain.MockAlgo{}

	sbPosition := 2
	bbPosition := 3
	tableID := RestCreatedTableFlopNext(t, sbPosition, bbPosition)

	restMakeAction(t, tableID.Hex(), domain.Check, bbPosition, secondPlayerToken)
	assertActionNewBettingRound(t, bbPosition, domain.Check, domain.FlopRound, bbPosition)
}

func TestTurn(t *testing.T) {
	t.Cleanup(cleanUp)
	domain.Algo = &domain.MockAlgo{}

	sbPosition := 0
	bbPosition := 1
	tableID := RestCreatedTableTurnNext(t, sbPosition, bbPosition)

	restMakeAction(t, tableID.Hex(), domain.Check, sbPosition, defaultToken)
	assertActionNewBettingRound(t, sbPosition, domain.Check, domain.TurnRound, bbPosition)
}

func TestRiver(t *testing.T) {
	t.Cleanup(cleanUp)
	domain.Algo = &domain.MockAlgo{}

	sbPosition := 2
	bbPosition := 3
	tableID := RestCreatedTableRiverNext(t, sbPosition, bbPosition)

	restMakeAction(t, tableID.Hex(), domain.Check, sbPosition, defaultToken)
	assertActionNewBettingRound(t, sbPosition, domain.Check, domain.RiverRound, bbPosition)
}
