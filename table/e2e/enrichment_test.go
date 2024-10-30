package e2e

import (
	"fmt"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/services/enrich"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestEnrichPlayers_ShouldSendSetPlayerConfig(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.Enrich.PlayersEnabled = true
	table := NewTableWithStartedGame(t, 0, 1, 0)
	insertTable(t, table)
	enrich.Players(table, table.AllPlayers())
	readMessage()
	msg := readMessage()
	fmt.Println(msg)
	assert.EqualValues(t, 2, len(msg.UserEvents.UserEvents))
	for _, e := range msg.UserEvents.UserEvents {
		assert.EqualValues(t, 1, len(e.Events))
		for _, event := range e.Events {
			assert.EqualValues(t, events.SetPlayerConfig, event.Type)
		}
	}
}

func TestEnrichPlayers_ShouldGetConfigFromDb(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.Enrich.PlayersEnabled = true
	table := NewTableWithStartedGame(t, 0, 1, 0)
	insertTable(t, table)

	body := `{"autoMuck":false, "autoTopUp": true, "autoRebuy": true}`
	rr := testRouter.Request(t, http.MethodPut, "/configs", &body, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)

	enrich.Players(table, table.AllPlayers())
	table = findTable(t, table.ID)
	firstConfig := table.GetPlayerUnsafe(0).AutoConfig
	assert.False(t, firstConfig.Muck)
	assert.True(t, firstConfig.TopUp)
	assert.True(t, firstConfig.ReBuy)
	secondConfig := table.GetPlayerUnsafe(1).AutoConfig
	assert.True(t, secondConfig.Muck)
	assert.False(t, secondConfig.TopUp)
	assert.False(t, secondConfig.ReBuy)
}
