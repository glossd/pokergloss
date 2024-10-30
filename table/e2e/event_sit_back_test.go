package e2e

import (
	"github.com/glossd/pokergloss/table/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSitBack(t *testing.T) {
	t.Cleanup(cleanUp)
	domain.Algo = &domain.MockAlgo{}

	table := NewTableWithStartedGame(t, 0, 1, -1)
	assert.Nil(t, table.MakeActionOnTimeout(0))
	assert.Nil(t, table.StartNextGame())
	insertTable(t, table)

	restSitBack(t, table.ID.Hex(), 0)
	assertSitBack(t, 0)

	// check if started a game
	assertStartHand(t, 1, 0, 0)
}

func assertSitBack(t *testing.T, position int) {
	assertMessage(t, 1, func(as []*Asserter) {
		as[0].assertSitBack(position)
	})
}
