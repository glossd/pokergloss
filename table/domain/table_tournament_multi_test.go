package domain

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMultiTable_AllInPlayer_ShouldNotDecide(t *testing.T) {
	lobby := startedMultiLobby(t, startedMultiParams{})
	assert.Len(t, lobby.GetTables(), 1)
	table := lobby.GetTables()[0]
	table.traise(t, 100)
	table.tcall(t)
	table.tfold(t)
	table.tfold(t)

	assert.Nil(t, table.StartNextGame())

	table.tallIn(t)
	table.tcall(t)
	table.tcall(t)

	assert.True(t, table.IsNewRound())

	table.tcheck(t)
	table.tallIn(t)
	table.tallIn(t)

	assert.True(t, table.IsGameEnd())
}

func TestMultiTable_2P_AllSittingOut(t *testing.T) {
	lobby := startedMultiLobby(t, startedMultiParams{userCount: 2})
	assert.Len(t, lobby.GetTables(), 1)
	table := lobby.GetTables()[0]

	table.ttimeout(t)
	assert.Nil(t, table.StartNextGame())
	table.ttimeout(t)
	assert.Nil(t, table.StartNextGame())

	assert.True(t, table.IsWaiting())
	assert.Len(t, table.AllPlayers(), 0)
}

func TestMultiTable_FeePercent(t *testing.T) {
	oldFee := conf.Props.Tournament.FeePercent
	conf.Props.Tournament.FeePercent = 0.02
	t.Cleanup(func() {
		conf.Props.Tournament.FeePercent = oldFee
	})
	Algo = Algo_2P_MockFirstPlayerLoses()
	lobby := startedMultiLobby(t, startedMultiParams{userCount: 2})
	assert.EqualValues(t, defaultBuyIn*2-10, lobby.PrizePool())
	assert.Len(t, lobby.GetTables(), 1)
	table := lobby.GetTables()[0]
	table.tallIn(t)
	table.tallIn(t)
	firstP := table.GetPlayerUnsafe(0)
	//secondP := table.GetPlayerUnsafe(1)
	assert.Nil(t, table.StartNextGame())
	assert.True(t, table.IsWaiting())
	assert.EqualValues(t, 0, firstP.GetTournamentInfo().Prize)
	// todo
	// The tournament info isn't set for multi in domain. The service fetches all the tables of tournament
	// and than calculates the place.
	//assert.EqualValues(t, defaultBuyIn*2-10, secondP.GetTournamentInfo().Prize)
}

func TestMultiTable_LongLobbyName(t *testing.T) {
	startedMultiLobby(t, startedMultiParams{name: "Super very long name of my awesome mulit lobby at some time"})
}

type startedMultiParams struct {
	// defaults to 3
	userCount int
	// defaults to 3
	tableSize int
	// defaults to "my multi"
	name string
}

func startedMultiLobby(t *testing.T, params startedMultiParams) *LobbyMulti {
	name := "my multi"
	count := 3
	tableSize := 3
	if params.userCount > 0 {
		count = params.userCount
	}
	if params.name != "" {
		name = params.name
	}
	if params.tableSize > 0 {
		tableSize = params.tableSize
	}
	lobby := NewLobbyMulti(NewLobbyMultiParams{
		StartAt:           timeutil.NowAdd(time.Millisecond),
		BuyIn:             250,
		BigBlind:          2,
		Name:              name,
		TableSize:         tableSize,
		LevelIncreaseTime: time.Minute,
		DecisionTimeout:   10 * time.Second,
	})
	for i := 0; i < count; i++ {
		assert.Nil(t, lobby.Register(getIden(i)))
	}
	lobby.Start()
	assert.NotNil(t, lobby.GetTables())
	return lobby
}
