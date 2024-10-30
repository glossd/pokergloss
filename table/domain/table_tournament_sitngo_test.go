package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSitAndGoTable_IncreaseBlinds(t *testing.T) {
	defaultMinutes := MinLevelIncreaseMinDuration
	t.Cleanup(func() {
		MinLevelIncreaseMinDuration = defaultMinutes
	})
	MinLevelIncreaseMinDuration = 0

	params := NewLobbySitAndGoParams{
		NewTableParams:    tableParams(2),
		BuyIn:             2500,
		PlacesPaid: 1,
		LevelIncreaseTime: 0,
	}

	table := fullLobbySitngo(t, params).GetTable()
	assert.EqualValues(t, 2, table.GetBigBlind().TotalRoundBet)
	assert.Nil(t, table.MakeAction(table.DecidingPosition, table.DecidingPlayerUnsafe().Identity, FoldAction))

	assert.True(t, table.IsGameEnd())
	assert.Nil(t, table.StartNextGame())

	assert.EqualValues(t, 4, table.BigBlind)
	assert.EqualValues(t, 4, table.GetBigBlind().TotalRoundBet)
}

func TestSitAndGoTable_SetPrizes(t *testing.T) {
	lobby := fullLobbySitngo(t, lobbySitngoParams(2, 1))
	assert.Len(t, lobby.Prizes, 1)
	assert.EqualValues(t, 1, lobby.Prizes[0].Place)
	assert.EqualValues(t, 500, lobby.Prizes[0].Prize)
}

func TestSitAndGoTable_StandPlayer(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}
	table := fullLobbySitngoSizes(t, 2, 1).GetTable()
	_, err := table.Stand(1, secondIdentity)
	assert.NotNil(t, err)
}

func TestSitAndGoTable_GivePrizesForLastPlayingPlayers(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = Algo_2P_MockSecondPlayerLoses()
	lobby := fullLobbySitngoSizes(t, 2, 1)
	table := lobby.GetTable()
	assert.Nil(t, table.MakeAction(0, firstIdentity, AllInAction))
	assert.Nil(t, table.MakeAction(1, secondIdentity, AllInAction))
	assert.True(t, table.IsGameEnd())
	assert.Nil(t, table.StartNextGame())

	assert.Len(t, table.NullifiedLeavingPlayers(), 2)
	assert.Len(t, table.TournamentWinners, 1)
	assert.EqualValues(t, defaultBuyIn * 2, table.TournamentWinners[0].Prize)
}

func TestSitAndGoTable_SittingOutPlayerAutoFolds_WhenHeTimeoutOnSmallBlind(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}
	lobby := fullLobbySitngoSizes(t, 2, 1)
	table := lobby.GetTable()

	assert.Nil(t, table.MakeActionOnTimeout(table.DecidingPosition))

	assert.True(t, table.IsGameEnd())

	assert.Nil(t, table.StartNextGame())

	table.traise(t, 20)

	assert.True(t, table.IsGameEnd())

	assert.Nil(t, table.StartNextGame())
}

func TestSitAndGoTable_SittingOutPlayerAutoFolds_WhenHeTimeoutsOnBigBlind(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}
	table := fullLobbySitngoSizes(t, 2, 1).GetTable()

	table.tcall(t)
	assert.Nil(t, table.MakeActionOnTimeout(table.DecidingPosition))

	assert.True(t, table.IsGameEnd())

	assert.Nil(t, table.StartNextGame())
	assert.True(t, table.IsGameEnd())

	assert.Nil(t, table.StartNextGame())

	table.traise(t, 20)
	assert.True(t, table.IsGameEnd())
}

func TestSitAndGoTable_AllIn_EndTable(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = Algo_2P_MockSecondPlayerLoses()
	table := fullLobbySitngoSizes(t, 2, 1).GetTable()

	table.tallIn(t)
	table.tallIn(t)

	assert.True(t, table.IsGameEnd())

	assert.Nil(t, table.StartNextGame())

	assert.Len(t, table.NullifiedLeavingPlayers(), 2)
	assert.True(t, table.IsWaiting())
}

func TestTestSitAndGoTable_3P_ResetOfSittingOutPlayer(t *testing.T) {
	table := fullLobbySitngoSizes(t, 3, 1).GetTable()

	table.tcall(t)
	table.tcall(t)
	table.tcheck(t)

	table.tcheck(t)
	table.tfold(t)
	sitOutPos := table.DecidingPosition
	assert.Nil(t, table.MakeActionOnTimeout(sitOutPos))

	assert.True(t, table.IsGameEnd())

	assert.Nil(t, table.StartNextGame())

	table.tcall(t)
	table.tcall(t)
	table.tsitBack(t, sitOutPos)
	assert.True(t, table.IsNewRound())

	table.tcheck(t)
	table.tfold(t)
	assert.True(t, table.IsGameEnd())

	assert.Nil(t, table.StartNextGame())

	table.tcall(t)
	table.tcall(t)
	table.tcheck(t)
	assert.True(t, table.IsFlop())
}

func TestSitAndGoTable_3P_FoldedSBPlayerNotActing(t *testing.T) {
	table := fullLobbySitngoSizes(t, 3, 1).GetTable()

	table.tcall(t)
	table.tcall(t)
	table.tfold(t)
	assert.True(t, table.IsNewRound())

	table.tcheck(t)
	table.tcheck(t)
	assert.True(t, table.IsNewRound())
}

func TestSitAndGoTable_2P_BlindAllIn(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)

	Algo = &MockAlgo{}
	table := fullLobbySitngoSizes(t, 2, 1).GetTable()
	table.traise(t, 248)
	table.tallIn(t)
	table.tfold(t)
	assert.True(t, table.IsGameEnd())
	assert.EqualValues(t, 1, table.GetPlayerUnsafe(0).Stack)

	Algo = Algo_2P_MockFirstPlayerLoses()
	assert.Nil(t, table.StartNextGame())
	assert.True(t, table.IsGameEnd())
	assert.Nil(t, table.StartNextGame())
	assert.True(t, table.IsWaiting())
}

func TestSitAndGoTable_2P_SbAllIn(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)

	Algo = &MockAlgo{}
	table := fullLobbySitngoSizes(t, 2, 1).GetTable()
	table.tcall(t)
	table.traise(t, 247)
	table.tallIn(t)
	table.tfold(t)
	assert.True(t, table.IsGameEnd())
	assert.EqualValues(t, 1, table.GetPlayerUnsafe(1).Stack)


	Algo = Algo_2P_MockSecondPlayerLoses()
	assert.Nil(t, table.StartNextGame())
	assert.True(t, table.IsGameEnd())
	assert.Nil(t, table.StartNextGame())
	assert.True(t, table.IsWaiting())
}

func TestSitAndGoTable_2P_AllSittingOut(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)

	Algo = &MockAlgo{}
	table := fullLobbySitngoSizes(t, 2, 1).GetTable()
	assert.Nil(t, table.MakeActionOnTimeout(0))
	assert.Nil(t, table.StartNextGame())
	assert.Nil(t, table.MakeActionOnTimeout(1))
	assert.Nil(t, table.StartNextGame())

	assert.True(t, table.IsWaiting())
	assert.Len(t, table.AllPlayers(), 0)
}
