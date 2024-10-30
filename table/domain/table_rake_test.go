package domain

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/stretchr/testify/assert"
	"testing"
)

func returnRakeParams(t *testing.T) {
	oldRP := conf.Props.Table.RakePercent
	oldMR := conf.Props.Table.MaxRake
	t.Cleanup(func() {
		conf.Props.Table.RakePercent = oldRP
		conf.Props.Table.MaxRake = oldMR
	})
}

func TestTableRake(t *testing.T) {
	returnRakeParams(t)
	conf.Props.Table.RakePercent = 0.01
	Algo = Algo_2P_MockSecondPlayerLoses()
	table := table2Player_gameRiver(t)
	table.tbet(t, 100)
	table.tcall(t)
	rake := 2
	assert.EqualValues(t, defaultBuyIn+100+defaultBB-rake, table.GetPlayerUnsafe(0).Stack)
	assert.EqualValues(t, defaultBuyIn-100-defaultBB, table.GetPlayerUnsafe(1).Stack)
}

func TestTableRake_NoRakeIfWinLessThanHundred(t *testing.T) {
	returnRakeParams(t)
	conf.Props.Table.RakePercent = 0.01
	Algo = Algo_2P_MockSecondPlayerLoses()
	table := table2Player_gameRiver(t)
	table.tcheck(t)
	table.tcheck(t)
	assert.EqualValues(t, defaultBuyIn+defaultBB, table.GetPlayerUnsafe(0).Stack)
	assert.EqualValues(t, defaultBuyIn-defaultBB, table.GetPlayerUnsafe(1).Stack)
}

func TestTableRake_DrawLowPots_NoRake(t *testing.T) {
	returnRakeParams(t)
	conf.Props.Table.RakePercent = 0.01
	Algo = Algo_2P_MockDraw()
	table := table2Player_gameRiver(t)
	table.tcheck(t)
	table.tcheck(t)
	assert.EqualValues(t, defaultBuyIn, table.GetPlayerUnsafe(0).Stack)
	assert.EqualValues(t, defaultBuyIn, table.GetPlayerUnsafe(1).Stack)
}

func TestTableRake_MaxRake(t *testing.T) {
	returnRakeParams(t)
	conf.Props.Table.RakePercent = 0.01
	conf.Props.Table.MaxRake = 2000
	Algo = Algo_2P_MockDraw()
	table, err := NewTable(NewTableParams{BigBlind: 1e8, Size: 6, Name: "Test", Identity: firstIdentity})
	assert.Nil(t, err)
	table.tsit(t, 0, 1e10)
	table.tsit(t, 1, 1e10)
	table.tallIn(t)
	table.tallIn(t)
	assert.EqualValues(t, 1e10-2000, table.GetPlayerUnsafe(0).Stack)
	assert.EqualValues(t, 1e10-2000, table.GetPlayerUnsafe(1).Stack)
}
