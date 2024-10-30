package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNL(t *testing.T) {
	Algo = &MockAlgo{}
	table := table2Players_bettingLimit(t, NL)

	table.tallIn(t)
	table.tallIn(t)
}

func TestPL(t *testing.T) {
	Algo = &MockAlgo{}
	table := table2Players_bettingLimit(t, PL)

	table.tallInFailed(t)
	table.traiseFailed(t, defaultBuyIn-50)
	table.tcall(t)

	table.tallInFailed(t)
	table.tbetFailed(t, defaultBuyIn-20)
	table.tcheck(t)

	table.tallInFailed(t)
	table.tbetFailed(t, defaultBB*2+1)
	table.tbet(t, defaultBB*2)
}

func TestML(t *testing.T) {
	Algo = &MockAlgo{}
	table := table2Players_bettingLimit(t, ML)

	table.tallInFailed(t)
	table.traiseFailed(t, defaultBuyIn-50)
	table.tcall(t)

	table.tallInFailed(t)
	table.tbetFailed(t, defaultBuyIn-20)
	table.tcheck(t)

	table.tbet(t, defaultBB*10)
	table.tallIn(t)
}

func TestFL(t *testing.T) {
	Algo = &MockAlgo{}
	table := table2Players_bettingLimit(t, FL)

	table.tallInFailed(t)
	table.traiseFailed(t, defaultBB*2)
	table.traise(t, defaultBB*1.5)

	table.tallInFailed(t)
	table.tcall(t)

	table.tbetFailed(t, defaultBB+1)
	table.tbet(t, defaultBB)
}

func table2Players_bettingLimit(t *testing.T, bl BettingLimit) *Table {
	params := getDefaultTableParams()
	params.BettingLimit = bl
	table, err := NewTable(params)
	assert.Nil(t, err)
	table.tsit(t, 0, defaultBuyIn)
	table.tsit(t, 1, defaultBuyIn)
	return table
}

func (t *Table) tallInFailed(test *testing.T) {
	t.tactionFailed(test, AllInAction)
}

func (t *Table) traiseFailed(test *testing.T, chips int64) {
	t.tactionFailed(test, RaiseAction(chips))
}

func (t *Table) tbetFailed(test *testing.T, chips int64) {
	t.tactionFailed(test, BetAction(chips))
}

func (t *Table) tactionFailed(test *testing.T, action Action) {
	assert.NotNil(test, t.MakeAction(t.DecidingPosition, getIden(t.DecidingPosition), action))
}
