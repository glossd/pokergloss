package e2e

import (
	"context"
	"fmt"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestAutoTopUp(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.Table.GameEndMinTimeout = -1
	domain.Algo = Algo_2P_MockZeroPositionLoses(t)
	tableID := RestCreatedTableEndNext(t, 0, 1)

	restTablePlayerAutoTop(t, tableID.Hex(), 0)

	restMakeAction(t, tableID.Hex(), domain.Check, 0)

	table := findTable(t, tableID)
	p := table.GetPlayerUnsafe(0)
	assert.EqualValues(t, defaultBuyIn-defaultBigBlind, p.Stack)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	actionhandler.DoStartGame(ctx, timeout.Key{TableID: tableID, Version: table.GameFlowVersion})

	table = findTable(t, tableID)
	p = table.GetPlayerUnsafe(0)
	assert.EqualValues(t, defaultBuyIn, p.Stack+p.TotalRoundBet)
}

func TestShouldNotTopUpOnSittingOut(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.MinDecisionTimeout = 0
	conf.Props.Table.GameEndMinTimeout = 0
	domain.Algo = &domain.MockAlgo{}

	table := InsertTableTimeout(t, 1)

	reserveAndBuyIn(t, 0, table, getToken(0))

	restTablePlayerAutoTop(t, table.ID.Hex(), 0)

	reserveAndBuyIn(t, 1, table, getToken(1))

	// timeout, gameEnd, stop
	table = findTable(t, table.ID)
	pOut := table.GetPlayerUnsafe(0)
	assert.EqualValues(t, domain.PlayerSittingOut, pOut.Status)
	assert.EqualValues(t, defaultBuyIn-defaultBigBlind/2, pOut.Stack)
}

func restTablePlayerAutoTop(t *testing.T, tableID string, pos int) {
	url := fmt.Sprintf("/tables/%s/seats/%d/configs/auto-top-up", tableID, pos)
	body := `{"autoTopUp": true}`
	rr := testRouter.Request(t, http.MethodPut, url, &body, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())
}
