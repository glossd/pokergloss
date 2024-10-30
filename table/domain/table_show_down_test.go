package domain

import (
	"github.com/stretchr/testify/assert"
"github.com/glossd/pokergloss/auth/authid"
"testing"
)

func TestShowDown(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = Algo_2P_MockFirstPlayerLoses()

	table := table2Player_gameRiver(t)
	table.tcheck(t)
	assert.Nil(t, table.SetAutoMuck(false, 0, getIden(0)))
	table.tcheck(t)
	assert.True(t, table.IsShowDown())

	assert.Nil(t, table.MakeShowDownAction(Show, 0, firstIdentity))
	assert.True(t, table.IsGameEnd())
}

func TestShowDown_EveryoneAllIn_onPreFlop_ShouldShowEveryone(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = Algo_2P_MockSecondPlayerLoses()

	table := table2Players_startedGame(t)
	assert.Nil(t, table.SetAutoMuck(false, 0, getIden(0)))
	assert.Nil(t, table.SetAutoMuck(false, 1, getIden(1)))
	table.tallIn(t)
	table.tallIn(t)
	assert.True(t, table.IsGameEnd())
}

func TestShowDownOnFold(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}

	table := table2Players_startedGame(t)
	assert.Nil(t, table.SetAutoMuck(false, 1, secondIdentity))
	assert.Nil(t, table.MakeAction(0, firstIdentity, FoldAction))
	assert.True(t, table.IsShowDown())

	assert.Nil(t, table.MakeShowDownAction(Show, 1, secondIdentity))

	assert.True(t, table.IsGameEnd())
}

func TestShowDown_WithOneSittingOutPlayer(t *testing.T) {
	Algo = &MockAlgo{}
	lobby := startedMultiLobby(t, startedMultiParams{tableSize: 6, userCount: 6})
	seats := lobby.IdensToSeats(lobby.Users)
	seats[0].Player.Stack = 49
	seats[1].Player = nil
	seats[2].Player = nil
	seats[3].Player.Stack = 120
	seats[4].Player.Stack = 707
	seats[5].Player = nil

	tableParams := NewTableParams{
		Name:            "1",
		Size:            lobby.TableSize,
		BigBlind:        2,
		DecisionTimeout: lobby.DecisionTimeout,
		Identity:        authid.Identity{UserId: "Automatic"},
	}

	Algo = NewMockFull(0,
		CardsStr(
			"Qd", "9c",
			"Kh", "4c",
		),
		CardsStr("7s", "Ah", "9s", "9h", "Jc"), // community cards
	)

	table, err := NewTableMulti(lobby, tableParams, seats)
	assert.Nil(t, err)
	assert.Nil(t, table.SetAutoMuck(false, 3, getIden(3)))
	assert.EqualValues(t, 0, table.DecidingPosition)
	table.tcall(t)
	assert.EqualValues(t, 3, table.DecidingPosition)
	table.tcall(t)
	assert.Nil(t, table.MakeActionOnTimeout(4))
	assert.True(t, table.IsNewRound())
	table.tcheck(t)
	table.tcheck(t)
	assert.True(t, table.IsNewRound())
	table.tcheck(t)
	table.tbet(t, 2)
	table.tcall(t)
	assert.True(t, table.IsNewRound())
	table.tcheck(t)
	table.tbet(t, 4)
	table.tcall(t)
	assert.True(t, table.IsShowDown())
	assert.EqualValues(t, 3, table.DecidingPosition)
	assert.EqualValues(t, 1, len(table.ShowedDownPlayers()))

	assert.Nil(t, table.MakeShowDownAction(Show, 3, getIden(3)))
	assert.True(t, table.IsGameEnd())
	assert.EqualValues(t, 1, len(table.ShowedDownPlayers()))
}

func TestShowdownTimeout(t *testing.T) {
	Algo = &MockAlgo{}
	lobby := startedMultiLobby(t, startedMultiParams{tableSize: 6, userCount: 6})
	seats := lobby.IdensToSeats(lobby.Users)
	seats[0].Player.Stack = 258 + 2
	seats[1].Player = nil
	seats[2].Player.Stack = 255
	seats[3].Player.Stack = 69
	seats[4].Player.Stack = 268 + 1
	seats[5].Player = nil

	tableParams := NewTableParams{
		Name:            "1",
		Size:            lobby.TableSize,
		BigBlind:        2,
		DecisionTimeout: lobby.DecisionTimeout,
		Identity:        authid.Identity{UserId: "Automatic"},
	}
	Algo = NewMockFull(3,
		CardsStr(
			"3d", "6s",
			"4d", "Td",
			"Kh", "Th",
			"Ts", "6h",
		),
		CardsStr("7s", "Kd", "7c", "3c"), // community cards
	)
	table, err := NewTableMulti(lobby, tableParams, seats)
	assert.Nil(t, err)

	assert.Nil(t, table.SetAutoMuck(false, 3, getIden(3)))

	assert.EqualValues(t, 2, table.DecidingPosition)
	table.tcall(t)
	table.tcall(t)
	table.tcall(t)
	table.tcheck(t)
	assert.True(t, table.IsNewRound())
	table.tcheck(t)
	table.tcheck(t)
	table.tcheck(t)
	table.tcheck(t)
	table.tcheck(t)
	table.tcheck(t)
	table.tcheck(t)
	assert.EqualValues(t, 3, table.DecidingPosition)
	table.tbet(t, 8)
	table.tfold(t)
	table.tfold(t)
	table.tfold(t)
	assert.Nil(t, table.MakeActionOnTimeout(table.DecidingPosition))
	assert.True(t, table.IsGameEnd())
}

func TestStackOverFlowPlayer_BeforeFold(t *testing.T) {
	table := table2Players_startedGame(t)
	table.traise(t, 200)
	assert.EqualValues(t, 203, table.TotalPot)
	table.tfold(t)
	assert.EqualValues(t, 4, table.TotalPot)
}

func TestStackOverFlowPlayerLast(t *testing.T) {
	table := table2Players_startedGame(t)
	table.tcall(t)
	assert.EqualValues(t, 4, table.TotalPot)
	table.traise(t, 200)
	assert.EqualValues(t, 204, table.TotalPot)
	table.tfold(t)
	assert.EqualValues(t, 4, table.TotalPot)
}

// https://youtu.be/w_WsO39YRx4?t=287
func TestStackOverFlow_OnFlop(t *testing.T) {
	Algo = &MockAlgo{}
	lobby := startedMultiLobby(t, startedMultiParams{tableSize: 6, userCount: 6})
	seats := lobby.IdensToSeats(lobby.Users)
	seats[0].Player.Stack = 1235
	seats[1].Player.Stack = 1769 + 200
	seats[2].Player.Stack = 387 + 400
	seats[3].Player.Stack = 574
	seats[4].Player = nil
	seats[5].Player = nil

	tableParams := NewTableParams{
		Name:            "1",
		Size:            lobby.TableSize,
		BigBlind:        400,
		DecisionTimeout: lobby.DecisionTimeout,
		Identity:        authid.Identity{UserId: "Automatic"},
	}
	Algo = NewMockFull(0,
		CardsStr(
			"2s", "7d",
			"Qh", "Ts",
			"8d", "Kh",
			"Js", "Qs",
		),
		CardsStr("5c", "8s", "Ks", "Td", "5s"), // community cards
	)
	table, err := NewTableMulti(lobby, tableParams, seats)
	assert.Nil(t, err)
	table.tcall(t)
	table.tfold(t)
	table.tcall(t)
	table.tcheck(t)
	assert.True(t, table.IsNewRound())
	table.tallIn(t)
	table.tallIn(t)
	table.tallIn(t)

	assert.EqualValues(t, 1182, table.GetPlayerUnsafe(1).Stack)
}

func TestShowdown_Bug_WhenMuckItShowed(t *testing.T) {
	Algo = NewMockFull(3,
		CardsStr(
			"Th", "Ks",
			"5h", "9c",
			"8h", "5c",
		),
		CardsStr("6h", "8c", "Qc", "7c", "2c"),
	)

	var table, err = NewTable(NewTableParams{Name: "table", BigBlind: 2, Size: 6, Identity: authid.Identity{UserId: "Automatic"}})
	assert.Nil(t, err)
	table.tsitUnsafe(2, defaultBuyIn)
	table.tsitUnsafe(3, defaultBuyIn)
	table.tsit(t, 4, defaultBuyIn)
	table.SetPlayerAutoConfig(table.GetPlayerUnsafe(2), PlayerAutoConfig{Muck: false})
	table.SetPlayerAutoConfig(table.GetPlayerUnsafe(3), PlayerAutoConfig{Muck: false})
	table.SetPlayerAutoConfig(table.GetPlayerUnsafe(4), PlayerAutoConfig{Muck: false})
	table.tcall(t)
	table.tcall(t)
	table.tcheck(t)

	table.tcheck(t)
	table.tcheck(t)
	table.tcheck(t)

	table.tcheck(t)
	table.tcheck(t)
	table.tcheck(t)

	table.tcheck(t)
	table.tcheck(t)
	table.tcheck(t)

	err = table.MakeShowDownAction(Muck, 2, getIden(2))
	assert.Nil(t, err)
	assert.EqualValues(t, Muck, table.GetPlayerUnsafe(2).ShowDownAction)
}

func TestShowdown_AskVipWhenLost(t *testing.T) {
	tableRiverForShowdownMarketCoins := func(t *testing.T) *Table {
		Algo = NewMockFull(2,
			CardsStr(
				"5h", "9c",
				"Ah", "Ks",
				"8h", "5c"),
			CardsStr("Ad", "Kc", "9c", "7h", "2s"))

		var table, err = NewTable(NewTableParams{Name: "table", BigBlind: 2, Size: 6, Identity: authid.Identity{UserId: "Automatic"}})
		assert.Nil(t, err)
		table.tsitUnsafe(1, defaultBuyIn)
		table.tsitUnsafe(2, defaultBuyIn)
		table.tsit(t, 3, defaultBuyIn)
		table.SetPlayerAutoConfig(table.GetPlayerUnsafe(1), PlayerAutoConfig{Muck: false})
		table.GetPlayerUnsafe(1).SetMarketItem("burger", 1)
		table.SetPlayerAutoConfig(table.GetPlayerUnsafe(2), PlayerAutoConfig{Muck: false})
		table.SetPlayerAutoConfig(table.GetPlayerUnsafe(3), PlayerAutoConfig{Muck: false})

		table.tcall(t)
		table.tcall(t)
		table.tcheck(t)

		table.tcheck(t)
		table.tbet(t, table.BigBlind)
		table.tcall(t)
		table.tfold(t)

		table.tcheck(t)
		table.tcheck(t)
		return table
	}

	t.Run("both check", func(t *testing.T) {
		table := tableRiverForShowdownMarketCoins(t)
		table.tcheck(t)
		table.tcheck(t)
		assert.Equal(t, 1, table.DecidingPosition)
		assert.Nil(t, table.MakeShowDownAction(Muck, 1, getIden(1)))
	})

	t.Run("aggression, who has coin item", func(t *testing.T) {
		table := tableRiverForShowdownMarketCoins(t)
		table.tbet(t, table.BigBlind)
		table.tcall(t)
		assert.Equal(t, 1, table.DecidingPosition)
		assert.Nil(t, table.MakeShowDownAction(Muck, 1, getIden(1)))
	})

	t.Run("auto muck", func(t *testing.T) {
		table := tableRiverForShowdownMarketCoins(t)
		table.SetPlayerAutoConfig(table.GetPlayerUnsafe(1), PlayerAutoConfig{Muck: true})
		table.tcheck(t)
		table.tcheck(t)
		assert.EqualValues(t, Muck, table.GetPlayerUnsafe(1).ShowDownAction)
		assert.EqualValues(t, Show, table.GetPlayerUnsafe(2).ShowDownAction)
	})
}
