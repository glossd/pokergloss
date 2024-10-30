package e2e

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/model"
	"github.com/glossd/pokergloss/table/web/client/mq"
	"sync"
	"testing"
)

func TestIntent_Put(t *testing.T) {
	t.Cleanup(cleanUp)

	tableID := RestCreatedTableWithStartedGame(t, 0, 1)

	intent := &model.Intent{Type: domain.CheckIntentType}
	restPutIntent(t, tableID.Hex(), 1, intent, secondPlayerToken)
	assertIntent(t, intent)

	restMakeAction(t, tableID.Hex(), domain.Call, 0)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	// todo doesn't work
	//go assertTwoActions(t, ws, wg)
	go assertTwoActions(t, wg)
	wg.Wait()
}

func assertTwoActions(t *testing.T, wg *sync.WaitGroup) {
	assertMessage(t, 4, func(as []*Asserter) {
		as[0].assertBettingAction(0, domain.Call)
		as[1].assertBettingAction(1, domain.Check)
		as[2].assertNewBettingRound(domain.FlopRound)
		as[3].assertTimeToDecide(1)
	})
	wg.Done()
}

func TestIntent_Delete(t *testing.T) {
	t.Cleanup(cleanUp)

	tableID := RestCreatedTableWithStartedGame(t, 0, 1)

	intent := &model.Intent{Type: domain.CheckIntentType}
	restPutIntent(t, tableID.Hex(), 1, intent, secondPlayerToken)
	assertIntent(t, intent)

	restDeleteIntent(t, tableID.Hex(), 1, secondPlayerToken)
	assertIntent(t, nil)
}

func TestIntent_Upgrade(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.Table.GameEndMinTimeout = 0

	domain.Algo = &domain.MockAlgo{}

	table := InsertTableTimeout(t, 0)
	hex := table.ID.Hex()

	reserveAndBuyIn(t, 0, table, getToken(0))
	reserveAndBuyIn(t, 1, table, getToken(1))
	reserveAndBuyInNoAsserts(t, 2, table, getToken(2))
	restMakeAction(t, hex, domain.Fold, 0)
	mq.ResetTestMQ()

	intent := &model.Intent{Type: domain.CheckFoldIntentType}
	restPutIntent(t, hex, 0, intent, getToken(0))
	assertIntent(t, intent)

	restMakeBetAction(t, hex, domain.Raise, 20, 1, getToken(1))
	assertSimpleAction(t, 1, 2, domain.Raise)

	assertIntent(t, model.ToIntent(&domain.FoldIntent))
}

func TestIntent_DeleteOnAction(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.Table.GameEndMinTimeout = 0

	domain.Algo = &domain.MockAlgo{}

	table := InsertTableTimeout(t, 0)
	hex := table.ID.Hex()

	reserveAndBuyIn(t, 0, table, getToken(0))
	reserveAndBuyIn(t, 1, table, getToken(1))
	reserveAndBuyInNoAsserts(t, 2, table, getToken(2))
	restMakeAction(t, hex, domain.Fold, 0)
	mq.ResetTestMQ()

	intent := &model.Intent{Type: domain.CheckIntentType}
	restPutIntent(t, hex, 0, intent, getToken(0))
	assertIntent(t, intent)

	restMakeBetAction(t, hex, domain.Raise, 20, 1, getToken(1))
	assertSimpleAction(t, 1, 2, domain.Raise)

	assertIntent(t, nil)
}

func assertIntent(t *testing.T, intent *model.Intent) {
	msg := readMessage()
	for _, events := range msg.UserEvents.UserEvents {
		for _, event := range events.Events {
			NewAsserter(t, event).assertIntent(intent)
		}
	}

}
