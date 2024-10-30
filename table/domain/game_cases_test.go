package domain

import (
	"github.com/stretchr/testify/assert"
"github.com/glossd/pokergloss/auth/authid"
"testing"
)

func Test_2P_PreFlop_SbFold(t *testing.T) {
	fPosition := 0
	sPosition := 1
	table := table2Players_startedGame(t)

	sb := table.DecidingPlayerUnsafe()
	table.tfold(t)

	assert.Len(t, table.Winners, 1)
	assert.EqualValues(t, 2, table.Winners[0].Chips)
	bbPosition := notPlayerPosition(sb, fPosition, sPosition)
	assert.Equal(t, bbPosition, table.Winners[0].Position)
	assert.Equal(t, "", table.Winners[0].HandRank)
}

func Test_2P_PreFlop_BbFold_After_SbRaise(t *testing.T) {
	table := table2Players_startedGame(t)
	table.traise(t, 10)
	table.tfold(t)

	assert.EqualValues(t, table.SmallBlindPosition(), table.Winners[0].Position)
	assert.EqualValues(t, 4, table.Winners[0].Chips)
}

func Test_2P_Flop_BbFold_After_SbBet(t *testing.T) {
	table := table2Players_startedGame(t)
	table.tcall(t)
	table.tcheck(t)
	table.tbet(t, 20)
	table.tfold(t)

	assert.True(t, table.IsGameEnd())
	assert.Len(t, table.Winners, 1)

	assert.EqualValues(t, table.BigBlindPosition(), table.Winners[0].Position)
	assert.EqualValues(t, 4, table.Winners[0].Chips)
}

func Test_3P_SecondPlayerTimeout_AfterPlayerJoined_ShouldKeepPlaying(t *testing.T) {
	table := table2PlayersPositions_startedGame(t, 0, 1)
	err := table.ReserveSeat(2, thirdIdentity)
	assert.Nil(t, err)
	_, err = table.BuyIn(defaultBuyIn, 2, thirdIdentity)
	assert.Nil(t, err)

	err = table.MakeActionOnTimeout(table.DecidingPosition)
	assert.Nil(t, err)

	assert.EqualValues(t, GameEndTable, table.Status)
}

func Test_2P_PlayerLostEverything(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = Algo_2P_MockSecondPlayerLoses()

	table := table2PlayersPositions_startedGame(t, 0, 1)

	err := table.MakeActionDeprecated(table.DecidingPosition, AllIn, 0)
	assert.Nil(t, err)
	err = table.MakeActionDeprecated(table.DecidingPosition, AllIn, 0)
	assert.Nil(t, err)

	assert.Len(t, table.Winners, 1)
	assert.EqualValues(t, 0, table.Winners[0].Position)

	assert.Zero(t, table.GetPlayerUnsafe(1).Stack)

	assert.True(t, table.IsGameEnd())

	assert.Nil(t, table.StartNextGame())
	assert.True(t, table.IsWaiting())
}

func TestCheckUserStackAfterWin(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{}
	lobby := startedMultiLobby(t, startedMultiParams{userCount: 6, tableSize: 6})
	assert.Len(t, lobby.GetTables(), 1)
	table := lobby.GetTables()[0]

	winnerPos := table.DecidingPosition
	table.traise(t, 20)
	table.tfold(t)
	table.tcall(t)
	table.ttimeout(t)
	table.tcall(t)
	table.tcall(t)
	assert.True(t, table.IsNewRound())

	table.tcheck(t)
	table.tcheck(t)
	table.tcheck(t)
	table.tcheck(t)
	assert.True(t, table.IsNewRound())

	table.tcheck(t)
	table.tcheck(t)
	table.tcheck(t)
	table.tcheck(t)
	assert.True(t, table.IsNewRound())

	table.tcheck(t)
	table.tbet(t, 3)
	table.traise(t, 100)
	table.tfold(t)
	table.tfold(t)
	table.tfold(t)
	assert.True(t, table.IsGameEnd())

	expectedStack := 250 + 63
	assert.EqualValues(t, expectedStack, table.GetPlayerUnsafe(winnerPos).Stack)
	assert.Nil(t, table.StartNextGame())
	assert.EqualValues(t, expectedStack-2, table.GetPlayerUnsafe(winnerPos).Stack)
}

// https://youtu.be/zfV02heMHGs?t=374
func Test_6P_NegativePot(t *testing.T) {
	lobby := startedMultiLobby(t, startedMultiParams{tableSize: 6, userCount: 6})
	seats := lobby.IdensToSeats(lobby.Users)
	seats[0].Player.Stack = 140
	seats[1].Player.Stack = 121
	seats[2].Player.Stack = 141
	seats[3].Player.Stack = 235 + 16
	seats[4].Player.Stack = 1001 + 32
	seats[5].Player.Stack = 64

	tableParams := NewTableParams{
		Name:            "1",
		Size:            lobby.TableSize,
		BigBlind:        32,
		DecisionTimeout: lobby.DecisionTimeout,
		Identity:        authid.Identity{UserId: "Automatic"},
	}
	Algo = NewMockFull(2,
		CardsStr(
			"Js", "Kh", // denis
			"8h", "8d", // eminem
			"3s", "Ad", // cheese
			"3d", "2c", // batman
			"5d", "4h", // robbin
			"7s", "8c", // vslv
		),
		CardsStr("Ac", "5c", "6c", "7h", "9s"), // community cards
	)
	table, err := NewTableMulti(lobby, tableParams, seats)
	assert.Nil(t, err)

	assert.EqualValues(t, 5, table.DecidingPosition)
	table.tcall(t)
	table.tcall(t)
	table.tcall(t)
	table.tcall(t)
	table.tcall(t)
	table.traise(t, 32)
	table.tallIn(t)
	table.tcall(t)
	table.tallIn(t)
	table.tcall(t)
	table.tcall(t)
	table.tcall(t)
	table.tallIn(t)
	table.tcall(t)
	table.tcall(t)
	table.tcall(t)
	assert.True(t, table.IsNewRound())
	assert.True(t, table.IsFlop())

	assert.Len(t, table.Pots, 4)
	assert.EqualValues(t, 64*6, table.Pots[0].Chips)        // 384
	assert.EqualValues(t, (121-64)*5, table.Pots[1].Chips)  // 285
	assert.EqualValues(t, (140-121)*4, table.Pots[2].Chips) // 76
	assert.EqualValues(t, 0, table.Pots[3].Chips)

	table.tcheck(t)
	table.tbet(t, 114)
	table.tallIn(t)
	table.tallIn(t)
	assert.True(t, table.IsGameEnd())

	assert.EqualValues(t, 4, table.LastAggressorPosition)

	assert.EqualValues(t, Show, table.GetPlayerUnsafe(4).ShowDownAction)
	assert.EqualValues(t, Show, table.GetPlayerUnsafe(5).ShowDownAction)
	assert.EqualValues(t, Show, table.GetPlayerUnsafe(0).ShowDownAction)
	assert.EqualValues(t, Show, table.GetPlayerUnsafe(1).ShowDownAction)
	assert.EqualValues(t, Show, table.GetPlayerUnsafe(2).ShowDownAction)
	assert.EqualValues(t, Show, table.GetPlayerUnsafe(3).ShowDownAction)

	assert.Len(t, table.Pots, 5)
	assert.EqualValues(t, 1*3, table.Pots[3].Chips)       // 3
	assert.EqualValues(t, (111-1)*2, table.Pots[4].Chips) // 220

	assert.EqualValues(t, 4, len(table.Winners))

	// todo get rid of damn domain.Table.Winners
	assert.EqualValues(t, table.Pots[0].Chips/2, table.Winners[1].Chips)
	assert.EqualValues(t, 192, table.Winners[1].Chips)
	assert.EqualValues(t, 5, table.Winners[1].Position)

	assert.EqualValues(t, table.Pots[0].Chips/2+table.Pots[1].Chips, table.Winners[0].Chips)
	assert.EqualValues(t, 192+285, table.Winners[0].Chips) // 477
	assert.EqualValues(t, 1, table.Winners[0].Position)

	assert.EqualValues(t, table.Pots[2].Chips+table.Pots[3].Chips, table.Winners[2].Chips)
	assert.EqualValues(t, 76+3, table.Winners[2].Chips)
	assert.EqualValues(t, 2, table.Winners[2].Position)

	assert.EqualValues(t, table.Pots[4].Chips, table.Winners[3].Chips)
	assert.EqualValues(t, 220, table.Winners[3].Chips)
	assert.EqualValues(t, 4, table.Winners[3].Position)
}

func TestForbidFoldActionAfterGameEnded(t *testing.T) {
	Algo = &MockAlgo{}
	lobby := startedMultiLobby(t, startedMultiParams{tableSize: 6, userCount: 6})
	seats := lobby.IdensToSeats(lobby.Users)
	seats[0].Player.Stack = 503
	seats[1].Player.Stack = 509
	seats[2].Player.Stack = 250
	seats[3].Player.Stack = 196
	seats[4].Player.Stack = 252
	seats[5].Player.Stack = 175

	tableParams := NewTableParams{
		Name:            "1",
		Size:            lobby.TableSize,
		BigBlind:        4,
		DecisionTimeout: lobby.DecisionTimeout,
		Identity:        authid.Identity{UserId: "Automatic"},
	}
	Algo = NewMockFull(5,
		CardsStr(
			"6s", "Ks",
			"3d", "7c",
			"2d", "2c",
			"Tc", "5h",
			"3c", "Qs",
			"9h", "Th",
		),
		CardsStr("2s", "9d", "8s", "Jd", "Ts"), // community cards
	)
	table, err := NewTableMulti(lobby, tableParams, seats)
	assert.Nil(t, err)

	assert.EqualValues(t, 2, table.DecidingPosition)
	table.tfold(t)
	table.tcall(t)
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
	assert.True(t, table.IsNewRound())

	table.tcheck(t)
	table.tcheck(t)
	table.tcheck(t)
	table.tcheck(t)
	table.tbet(t, 4)
	table.tcall(t)
	table.tfold(t)
	table.tfold(t)
	table.tfold(t)
	assert.True(t, table.IsNewRound())

	table.tbet(t, 383)
	lastDecidingPosition := table.DecidingPosition
	table.tallIn(t)
	assert.True(t, table.IsGameEnd())

	assert.NotNil(t, table.MakeAction(lastDecidingPosition, getIden(lastDecidingPosition), FoldAction))
}

// https://youtu.be/MWUAkmvULmE?t=204
func TestWrongSecondPot_ShouldBeOnlyOne(t *testing.T) {
	Algo = &MockAlgo{}
	lobby := startedMultiLobby(t, startedMultiParams{tableSize: 6, userCount: 6})
	seats := lobby.IdensToSeats(lobby.Users)
	seats[0].Player.Stack = 886
	seats[1].Player.Stack = 1876 + 50
	seats[1].Player.Status = PlayerSittingOut
	seats[2].Player.Stack = 341
	seats[3].Player.Stack = 886
	seats[4].Player.Stack = 711
	seats[5].Player.Stack = 513

	tableParams := NewTableParams{
		Name:            "1",
		Size:            lobby.TableSize,
		BigBlind:        100,
		DecisionTimeout: lobby.DecisionTimeout,
		Identity:        authid.Identity{UserId: "Automatic"},
	}
	Algo = NewMockFull(0,
		CardsStr(
			"Kh", "7s",
			"Ks", "Qh",
			"Td", "As",
			"9d", "5c",
			"9s", "Ah",
			"5h", "7h",
		),
		CardsStr("Kc", "Ad", "Qc", "7d", "Th"), // community cards
	)
	table, err := NewTableMulti(lobby, tableParams, seats)
	assert.Nil(t, err)
	assert.EqualValues(t, 3, table.DecidingPosition)

	table.tfold(t)
	table.tcall(t)
	table.tfold(t)
	table.tfold(t)
	assert.EqualValues(t, 2, len(table.MadeActionPlayerPositions()))
	table.tallIn(t)
	table.tcall(t)
	assert.True(t, table.IsGameEnd())
	assert.EqualValues(t, 341*2+50, table.Pots[0].Chips)
}

func TestStackOverFlowStack_onShowdown(t *testing.T) {

	Algo = NewMockFull(1,
		CardsStr(
			"Ac", "Th",
			"4s", "7s",
		),
		CardsStr("Jd", "Qh", "8d", "Kh", "2c"), // community cards
	)

	table := table2Players_startedGame(t)
	table.GetPlayerUnsafe(0).Stack = 98
	table.GetPlayerUnsafe(1).Stack = 43
	assert.Nil(t, table.SetAutoMuck(false, 1, getIden(1)))

	table.tcall(t)
	table.tcheck(t)

	table.tcheck(t)
	table.tcheck(t)

	table.tcheck(t)
	table.tcheck(t)

	table.tallIn(t)
	table.tallIn(t)

	sofp := table.BuildStackOverflowPlayer()
	assert.EqualValues(t, 100-44, sofp.Stack)
}

func TestPotNegative1(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = NewMockFull(3,
		CardsStr(
			"4c", "3h",
			"5c", "5h",
			"Td", "3c",
		),
		CardsStr("4d", "6s", "7d", "4h", "Ah"), // community cards
	)

	var table, err = NewTable(NewTableParams{Name: "table", BigBlind: 2, Size: 9, Identity: authid.Identity{UserId: "Automatic"}})
	assert.Nil(t, err)
	table.tsitUnsafe(2, 1065)
	table.tsitUnsafe(3, 1585)
	table.tsitUnsafe(6, 0)
	table.Seats[6].Player.Status = PlayerSittingOut
	table.tsit(t, 4, 400)

	assert.EqualValues(t, 3, table.DecidingPosition)
	table.tcall(t)
	assert.Nil(t, table.MakeActionOnTimeout(4))
	table.tcheck(t)

	table.tallIn(t)
	table.tcall(t)

	assert.EqualValues(t, 1, len(table.Pots))
	assert.EqualValues(t, 1065*2+1, table.Pots[0].Chips)
}

func TestPots_PotLost(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = NewMockFull(2,
		CardsStr(
			"Jd", "4s",
			"Kh", "5s",
			"2s", "3d",
			"9d", "Jc",
			"7s", "Tc",
			"Jh", "9s",
		),
		CardsStr("9c", "2c", "3s", "3d", "Js"), // community cards
	)

	table, err := NewTable(NewTableParams{Name: "table", BigBlind: 2, Size: 6, Identity: authid.Identity{UserId: "Automatic"}})
	assert.Nil(t, err)
	table.tsitUnsafe(0, 84)
	table.tsitUnsafe(1, 522+1977)
	table.tsitUnsafe(2, 126)
	table.tsitUnsafe(3, 476)
	table.tsitUnsafe(4, 522)
	table.tsit(t, 5, 122)

	assert.EqualValues(t, 5, table.DecidingPosition)
	table.tallIn(t)
	table.tallIn(t)
	table.tcall(t)
	table.tallIn(t)
	table.tallIn(t)
	table.tcall(t)
	table.tcall(t)

	table.tcheck(t)
	table.tcheck(t)

	table.tcheck(t)
	table.tcheck(t)

	table.tcheck(t)
	table.tbet(t, 73)
	table.tallIn(t)

	assert.EqualValues(t, 5, len(table.Pots))

	assertOnePotWin(t, table, 0, 2, 504)
	assertOnePotWin(t, table, 1, 2, 190)
	assertOnePotWin(t, table, 2, 2, 16)
	assertOnePotWin(t, table, 3, 3, 1050)
	assertOnePotWin(t, table, 4, 1, 92)
}

func TestSmallBlindStackLessThanBB(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{dealerPos: 1}
	table := table2Players_startedGame(t)
	table.traise(t, defaultBuyIn-2)
	table.tcall(t)
	table.tfold(t)
	assert.EqualValues(t, 1, table.GetPlayerUnsafe(0).Stack)
	assert.EqualValues(t, defaultBuyIn*2-1, table.GetPlayerUnsafe(1).Stack)
	Algo = Algo_2P_MockFirstPlayerLoses()
	assert.Nil(t, table.StartNextGame())
	assert.True(t, table.IsRiver())
}

func TestPlayerWithStackLessThanBigBlindWins(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{dealerPos: 1}
	table, err := NewTable(NewTableParams{
		Name:            "my table",
		Size:            defaultTableSize,
		BigBlind:        10,
		BettingLimit:    NL,
		DecisionTimeout: 0,
		Identity:        firstIdentity,
	})
	assert.Nil(t, err)
	table.tsit(t, 0, 500)
	table.tsit(t, 1, 500)
	table.traise(t, 500-5-4)
	table.tcall(t)
	table.tcheck(t)
	table.tfold(t)
	assert.EqualValues(t, 4, table.GetPlayerUnsafe(1).Stack)
	Algo = Algo_2P_MockFirstPlayerLoses()
	assert.Nil(t, table.StartNextGame())
	assert.True(t, table.IsRiver())
	assert.EqualValues(t, 8, table.GetPlayerUnsafe(1).Stack)
	Algo = Algo_2P_MockSecondPlayerLoses()
	assert.Nil(t, table.StartNextGame())
	assert.EqualValues(t, DealerSmallBlind, table.GetPlayerUnsafe(1).Blind)
	assert.True(t, table.IsPreFlop())
}

func TestPlayerWithStackLessThanBigBlindWinsAndGetToBeBigBlind(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = &MockAlgo{dealerPos: 1}
	table := table2Players_startedGame(t)
	table.traise(t, defaultBuyIn-2)
	table.tcall(t)
	table.tfold(t)
	assert.EqualValues(t, 1, table.GetPlayerUnsafe(0).Stack)
	Algo = Algo_2P_MockSecondPlayerLoses()
	assert.Nil(t, table.StartNextGame())
	assert.True(t, table.IsRiver())
	assert.EqualValues(t, 2, table.GetPlayerUnsafe(0).Stack)
	Algo = Algo_2P_MockFirstPlayerLoses()
	assert.Nil(t, table.StartNextGame())
	assert.EqualValues(t, BigBlind, table.GetPlayerUnsafe(0).Blind)
	assert.True(t, table.IsPreFlop())
	table.tcall(t)
	assert.True(t, table.IsRiver())
}

func TestBrokePlayerToReservedSeat(t *testing.T) {
	t.Cleanup(cleanToRealAlgo)
	Algo = Algo_2P_MockFirstPlayerLoses()
	table := table2Players_startedGame(t)
	table.tallIn(t)
	table.tallIn(t)
	assert.Nil(t, table.StartNextGame())
	assert.True(t, table.IsWaiting())
	assert.EqualValues(t, PlayerReservedSeat, table.GetPlayerUnsafe(0).Status)
}

func assertOnePotWin(t *testing.T, table *Table, potIdx, winPos int, chips int64) {
	assert.EqualValues(t, 1, len(table.Pots[potIdx].WinnerPositions))
	assert.EqualValues(t, winPos, table.Pots[potIdx].WinnerPositions[0])
	assert.EqualValues(t, chips, table.Pots[potIdx].Chips)
}

func (t *Table) tsitUnsafe(pos int, chips int64) {
	t.Seats[pos].addPlayer(getIden(pos))
	t.Seats[pos].Player.setInitStack(chips)
}

func assertPotWinnerPosition(t *testing.T, pot *Pot, pos int) {
	assert.EqualValues(t, 1, len(pot.WinnerPositions))
	assert.EqualValues(t, pos, pot.WinnerPositions[0])
}
