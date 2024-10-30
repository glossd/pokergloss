package e2e

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/services/cleaning"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCleanAlonePlayerOnPersistentTable(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.Cleaning.AlonePlayerOnPersistentTable = 0
	conf.Props.Table.GameEndMinTimeout = 0
	table := NewTableTimeout(t, -1)
	table.IsPersistent = true
	insertTable(t, table)
	reserveAndBuyIn(t, 0, table, defaultToken)

	assert.EqualValues(t, 1, cleaning.CleanAlonePlayerOnPersistentTable())
	assert.EqualValues(t, 0, len(findTable(t, table.ID).AllPlayers()))
}
