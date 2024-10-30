package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSecondPlayerSitBack_ShouldStartGame(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}
	table := table2PlayersPositions_startedGame(t, 0, 1)
	assert.Nil(t, table.MakeActionOnTimeout(0))
	assert.Nil(t, table.StartNextGame())

	newGame, err := table.SitBack(0, firstIdentity)
	assert.Nil(t, err)

	assert.True(t, newGame)

	assert.EqualValues(t, PlayingTable, table.Status)
}
