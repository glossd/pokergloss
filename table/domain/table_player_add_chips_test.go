package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTable_AddChips_ToStack(t *testing.T) {
	table := defaultTable(t)
	assert.Nil(t, table.ReserveSeat(0, firstIdentity))
	_, err := table.BuyIn(defaultBuyIn, 0, firstIdentity)
	assert.Nil(t, err)
	assert.Nil(t, table.AddChips(10, 0, firstIdentity))
	assert.EqualValues(t, defaultBuyIn+ 10, table.GetPlayerUnsafe(0).Stack)
}

func TestTable_AddChips_OnGameStart(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = Algo_2P_MockSecondPlayerLoses()

	table := table2PlayersPositions_startedGame(t, 0, 1)
	assert.Nil(t, table.AddChips(10, 0, firstIdentity))
	assert.Nil(t, table.MakeActionDeprecated(0, Fold, 0))

	assert.Nil(t, table.StartNextGame())

	assert.EqualValues(t, 257, table.GetPlayerUnsafe(0).Stack)
}
