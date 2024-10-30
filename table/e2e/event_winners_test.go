package e2e

import (
	"github.com/glossd/pokergloss/table/domain"
	"testing"
)

func TestWinners(t *testing.T) {
	t.Cleanup(cleanUp)
	domain.Algo = &domain.MockAlgo{}

	sbPosition := 0 // default
	bbPosition := 1 // second
	tableID := RestCreatedTableEndNext(t, sbPosition, bbPosition)

	restMakeAction(t, tableID.Hex(), domain.Check, sbPosition)
	assertActionGameEndWithShowdown(t, sbPosition, domain.Check, bbPosition, sbPosition, 2)
}

func assertActionGameEndWithShowdown(t *testing.T, pos int, action domain.ActionType, bbPosition int, sbPosition int, winnersCount int) {
	assertMessage(t, 4, func(as []*Asserter) {
		as[0].assertBettingAction(pos, action)
		as[1].assertShowdown(bbPosition, false)
		as[2].assertShowdown(sbPosition, false)
		as[3].assertWinners(winnersCount)
	})
}

func assertActionGameEndWithOneShowdown(t *testing.T, pos int, winPosition int) {
	assertMessage(t, 3, func(as []*Asserter) {
		as[0].assertBettingAction(pos, domain.Fold)
		as[1].assertShowdown(winPosition, true)
		as[2].assertWinners(1)
	})
}

func TestGameEndFold(t *testing.T) {
	t.Cleanup(cleanUp)

	sbPosition := 0 // default
	bbPosition := 1 // second
	tableID := RestCreatedTableFlopNext(t, sbPosition, bbPosition)

	restMakeAction(t, tableID.Hex(), domain.Fold, bbPosition, secondPlayerToken)
	assertActionFoldAndWinners(t, bbPosition, sbPosition)
}

func assertActionFoldAndWinners(t *testing.T, pos int, winnerPosition int) {
	assertMessage(t, 3, func(as []*Asserter) {
		as[0].assertBettingAction(pos, domain.Fold)
		as[1].assertShowdown(winnerPosition, true)
		as[2].assertWinners(1)
	})
}
