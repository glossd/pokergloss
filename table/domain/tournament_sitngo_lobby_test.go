package domain

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLobby(t *testing.T) {
	l, err := NewLobbySitAndGo(lobbySitngoParams(2, 1))
	assert.Nil(t, err)

	err = l.Register(firstIdentity, 0)
	assert.Nil(t, err)
	err = l.Register(secondIdentity, 1)
	assert.Nil(t, err)
	assert.NotNil(t, l.GetTable())

	assert.Len(t, l.GetTable().DecidablePlayers(), 2)
	assert.EqualValues(t, 1, len(l.Prizes))
	assert.EqualValues(t, 2*defaultBuyIn, l.Prizes[0].Prize)
}

func TestLobbySitAndGo_ComputeTournamentPrize(t *testing.T) {
	for i := 2; i < 8; i++ {
		assertOnePaidPlace(t, fullLobbySitngoSizes(t, i, 1))
	}

	lobby := fullLobbySitngoSizes(t, 6, 2)
	assert.EqualValues(t, 250*6*3/4, ComputeTournamentPrize(lobby, 1))
	assert.EqualValues(t, 250*6*1/4, ComputeTournamentPrize(lobby, 2))
	assert.EqualValues(t, 0, ComputeTournamentPrize(lobby, 3))

	lobby = fullLobbySitngoSizes(t, 9, 3)
	assert.EqualValues(t, 250*9*3/6, ComputeTournamentPrize(lobby, 1))
	assert.EqualValues(t, 250*9*2/6, ComputeTournamentPrize(lobby, 2))
	assert.EqualValues(t, 250*9*1/6, ComputeTournamentPrize(lobby, 3))
	assert.EqualValues(t, 0, ComputeTournamentPrize(lobby, 4))
}

func TestLobbySitAndGo_BackToGame(t *testing.T) {

	lobby := fullLobbySitngoSizes(t, 2, 1)
	table := lobby.GetTable()
	err := table.MakeActionOnTimeout(table.DecidingPosition)
	assert.Nil(t, err)

	assert.True(t, table.IsGameEnd())

	assert.Nil(t, table.StartNextGame())
	table.tcall(t)
	assert.True(t, table.IsGameEnd())

	assert.Nil(t, table.StartNextGame())
	assert.True(t, table.IsGameEnd())

	assert.Nil(t, table.StartNextGame())

	nextPos := table.nextDecidingPos_2P()
	table.tsitBack(t, nextPos)

	table.tcall(t)
	assert.True(t, table.IsGameEnd())

	assert.Nil(t, table.StartNextGame())

	table.tcall(t)
	table.tcheck(t)
	assert.True(t, table.IsNewRound())
}

func TestLobbySitAndGo_StartByStartAt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		lobby, err := NewLobbySitAndGo(lobbySitngoParams(3, 1))
		assert.Nil(t, err)
		assert.Nil(t, lobby.Register(getIden(0), 0))
		assert.Nil(t, lobby.Register(getIden(1), 1))

		assert.Nil(t, lobby.StartAnyway())
		assert.EqualValues(t, 3, lobby.GetTable().Size)
		assert.EqualValues(t, 2, len(lobby.GetTable().AllPlayers()))
		assert.EqualValues(t, lobby.BuyIn*2, lobby.PrizePool())
	})
}

func TestLobbySitNAndGo_FeePercent(t *testing.T) {
	prevFee := conf.Props.Tournament.FeePercent
	conf.Props.Tournament.FeePercent = 0.02
	t.Cleanup(func() {
		conf.Props.Tournament.FeePercent = prevFee
	})
	Algo = Algo_2P_MockFirstPlayerLoses()
	lobby, err := NewLobbySitAndGo(lobbySitngoParams(3, 1))
	assert.Nil(t, err)
	assert.Nil(t, lobby.Register(getIden(0), 0))
	assert.Nil(t, lobby.Register(getIden(1), 1))
	assert.Nil(t, lobby.StartAnyway())
	assert.EqualValues(t, 500-10, lobby.PrizePool())
	assert.EqualValues(t, 1, lobby.PlacesPaid)

	table := lobby.GetTable()
	table.tallIn(t)
	table.tallIn(t)
	firstP := table.GetPlayerUnsafe(0)
	secondP := table.GetPlayerUnsafe(1)
	assert.Nil(t, table.StartNextGame())
	assert.EqualValues(t, 0, firstP.GetTournamentInfo().Prize)
	assert.EqualValues(t, 500-10, secondP.GetTournamentInfo().Prize)
}

func assertOnePaidPlace(t *testing.T, lobby *LobbySitAndGo) {
	assert.EqualValues(t, defaultBuyIn*lobby.Size, ComputeTournamentPrize(lobby, 1))
	for i := 2; i <= lobby.Size; i++ {
		assert.EqualValues(t, 0, ComputeTournamentPrize(lobby, i))
	}
}

func fullLobbySitngo(t *testing.T, params NewLobbySitAndGoParams) *LobbySitAndGo {
	lobby, err := NewLobbySitAndGo(params)
	assert.Nil(t, err)
	for i := 0; i < lobby.Size; i++ {
		err := lobby.Register(getIden(i), i)
		assert.Nil(t, err)
		if i == lobby.Size {
			assert.NotNil(t, lobby.GetTable())
		}
	}
	assert.Nil(t, err)
	return lobby
}

func fullLobbySitngoSizes(t *testing.T, size, paidPlaces int) *LobbySitAndGo {
	lobby, err := NewLobbySitAndGo(lobbySitngoParams(size, paidPlaces))
	assert.Nil(t, err)
	for i := 0; i < lobby.Size; i++ {
		err := lobby.Register(getIden(i), i)
		assert.Nil(t, err)
		if i == lobby.Size {
			assert.NotNil(t, lobby.GetTable())
		}
	}
	assert.Nil(t, err)
	return lobby
}

func lobbySitngoParams(size, placesPaid int) NewLobbySitAndGoParams {
	params := NewLobbySitAndGoParams{
		NewTableParams:    tableParams(size),
		PlacesPaid:        placesPaid,
		BuyIn:             defaultBuyIn,
		LevelIncreaseTime: time.Minute,
	}
	return params
}
