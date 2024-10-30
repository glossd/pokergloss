package e2e

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/cleaning"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCleanWaitingTables(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.Cleaning.WaitingTablesTimeout = 0
	conf.Props.Table.GameEndMinTimeout = 0
	tableID := RestCreatedTableWithStartedGame(t, 0, 1)
	assert.False(t, actionhandler.DoDecisionTimeoutNoCtx(timeout.Key{TableID: tableID, Position: 0, Version: 2}))
	table := findTable(t, tableID)
	assert.EqualValues(t, domain.WaitingTable, table.Status)
	assert.EqualValues(t, 1, cleaning.CleanWaitingTables())
}
