package e2e

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSitBackAndStand_WhileGameEnd(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.Table.GameEndMinTimeout = -1
	domain.Algo = &domain.MockAlgo{}

	tableID := RestCreatedTableWithStartedGame(t, 0, 1)
	actionhandler.DoDecisionTimeoutNoCtx(timeout.Key{TableID: tableID, Position: 0, Version: 2})
	table := findTable(t, tableID)
	restSitBack(t, tableID.Hex(), 0, getToken(0))
	restPlayerStand(t, tableID, 0, getToken(0))
	tableNew := findTable(t, tableID)
	assert.EqualValues(t, table.GameFlowVersion, tableNew.GameFlowVersion)
}
