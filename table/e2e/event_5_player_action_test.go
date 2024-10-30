package e2e

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"testing"
	"time"
)

func TestPlayerAction(t *testing.T) {
	t.Cleanup(cleanUp)
	domain.Algo = &domain.MockAlgo{}

	sbPosition := 0
	bbPosition := 1
	tableId := RestCreatedTableWithStartedGame(t, sbPosition, bbPosition)

	restMakeAction(t, tableId.Hex(), domain.Call, sbPosition, defaultToken)
	assertSimpleAction(t, sbPosition, bbPosition, domain.Call)
}

func TestPlayerAction_BeforeWaitDuration(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.PlayerActionDuration = time.Second
	domain.Algo = &domain.MockAlgo{}

	sbPosition := 0
	bbPosition := 1
	tableId := RestCreatedTableWithStartedGame(t, sbPosition, bbPosition)

	restMakeAction(t, tableId.Hex(), domain.Call, sbPosition, getToken(sbPosition))
	restMakeAction(t, tableId.Hex(), domain.Check, bbPosition, getToken(bbPosition))
	rr := restMakeActionStatus(t, tableId.Hex(), domain.Check, bbPosition, 400, getToken(bbPosition))
	assert.EqualValues(t, "wait...", gjson.Get(rr.Body.String(), "message").String())
}

func TestPlayerAction_ChipsBeforeGameEndResult(t *testing.T) {
	prevPropsSetup(t)
	domain.Algo = Algo_2P_MockZeroPositionLoses(t)
	tableId := RestCreatedTableWithStartedGame(t, 0, 1)

	restMakeAction(t, tableId.Hex(), domain.AllIn, 0, defaultToken)
	assertSimpleAction(t, 0, 1, domain.AllIn)
	restMakeAction(t, tableId.Hex(), domain.AllIn, 1, secondPlayerToken)
	assertMessage(t, 4, func(as []*Asserter) {
		as[0].assertBettingAction(1, domain.AllIn)
		assert.Zero(t, gjson.Get(as[0].Event.Payload, "table.seats.0.player.stack").Int())
		assert.EqualValues(t, defaultBuyIn, gjson.Get(as[0].Event.Payload, "table.seats.0.player.totalRoundBet").Int())
	})
}

func assertSimpleAction(t *testing.T, pos, decidingPos int, action domain.ActionType) {
	assertMessage(t, 2, func(as []*Asserter) {
		as[0].assertBettingAction(pos, action)
		as[1].assertTimeToDecide(decidingPos)
	})
}
