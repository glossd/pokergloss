package e2e

import (
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/web/client/mq"
	"github.com/glossd/pokergloss/table/web/client/mqpub"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGameEndEvent_AllFoldOnPreFlop_ShouldSetWinnerRank(t *testing.T) {
	t.Cleanup(cleanUp)
	tableID := RestCreatedTableTurnNext(t, 0, 1)
	table := findTable(t, tableID)
	assert.Nil(t, table.MakeAction(0, getIden(0), domain.FoldAction))

	assert.EqualValues(t, 1, len(table.Winners))
	assert.EqualValues(t, "", table.Winners[0].HandRank)

	mqpub.PublishGameEnd(table)
	msg := <-mq.TestGameEndMQ
	assert.EqualValues(t, 1, len(msg.Winners))
	assert.EqualValues(t, "Full House", msg.Winners[0].Hand)
}
