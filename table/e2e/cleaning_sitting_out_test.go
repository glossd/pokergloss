package e2e

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/services/cleaning"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestCleaningSittingOutPlayer(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.Cleaning.CashSittingOutTimeout = 0

	var tableIds []primitive.ObjectID
	for i := 0; i < 2; i++ {
		table := NewTableWithStartedGame(t, 0, 1, -1)
		assert.Nil(t, table.MakeActionOnTimeout(0)) // player sits out
		insertTable(t, table)
		tableIds = append(tableIds, table.ID)
		assert.Len(t, table.AllPlayers(), 2)
	}

	cleaning.CleanSittingOutPlayers()

	for _, id := range tableIds {
		table, err := db.FindTableNoCtx(id)
		assert.Nil(t, err)
		assert.Len(t, table.AllPlayers(), 1)
	}

	assertPlayerLeft(t, 0)
	assertPlayerLeft(t, 0)
}
