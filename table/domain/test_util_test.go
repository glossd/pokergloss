package domain

import (
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

var (
	firstIdentity   = authid.Identity{Username: "1", UserId: "1ID"}
	secondIdentity  = authid.Identity{Username: "2", UserId: "2ID"}
	thirdIdentity   = authid.Identity{Username: "3", UserId: "3ID"}
	fourthIdentity  = authid.Identity{Username: "4", UserId: "4ID"}
	fifthIdentity   = authid.Identity{Username: "5", UserId: "5ID"}
	sixthIdentity   = authid.Identity{Username: "6", UserId: "6ID"}
	seventhIdentity = authid.Identity{Username: "7", UserId: "7ID"}
	eighthIdentity  = authid.Identity{Username: "8", UserId: "8ID"}
	ninthIdentity   = authid.Identity{Username: "9", UserId: "9ID"}
)

const (
	defaultBuyIn     = 250
	defaultBB        = 2
	defaultTableSize = 9
)

var defaultTableParams = NewTableParams{
	Name:            "my table",
	Size:            defaultTableSize,
	BigBlind:        defaultBB,
	BettingLimit:    NL,
	DecisionTimeout: 0,
	Identity:        firstIdentity,
}

var cleanToRealAlgo = func() { Algo = &RealAlgo{} }

func getIden(position int) authid.Identity {
	var i authid.Identity
	switch position {
	case 0:
		i = firstIdentity
	case 1:
		i = secondIdentity
	case 2:
		i = thirdIdentity
	case 3:
		i = fourthIdentity
	case 4:
		i = fifthIdentity
	case 5:
		i = sixthIdentity
	case 6:
		i = seventhIdentity
	case 7:
		i = eighthIdentity
	case 8:
		i = ninthIdentity
	}
	return i
}

func (t *Table) nextDecidingPos_2P() int {
	return int(math.Abs(float64(t.DecidingPosition - 1)))
}

func getDefaultTableParams() NewTableParams {
	return defaultTableParams
}

// Deprecated
// Returns player with small blind
func table2Players_startedGameDeprecated(t *testing.T, fPosition int, sPosition int) (*Table, *Player) {
	table := table2PlayersPositions_startedGame(t, fPosition, sPosition)
	return table, table.DecidingPlayerUnsafe()
}

func table2PlayersPositions_startedGame(t *testing.T, fPosition int, sPosition int) *Table {
	return table2Players_startedGame_chips(t, fPosition, sPosition, defaultBuyIn, defaultBuyIn)
}

func table2Players_startedGame(t *testing.T) *Table {
	return table2Players_startedGame_chips(t, 0, 1, defaultBuyIn, defaultBuyIn)
}

func table2Players_flop(t *testing.T) *Table {
	table := table2Players_startedGame(t)
	table.tcall(t)
	table.tcheck(t)
	return table
}

func table2Players_startedGame_chips(t *testing.T, fPosition, sPosition int, fChips, sChips int64) *Table {
	table := defaultTable(t)
	err := table.ReserveSeat(fPosition, firstIdentity)
	assert.Nil(t, err)
	isNewGame, err := table.BuyIn(fChips, fPosition, firstIdentity)
	assert.Nil(t, err)
	assert.False(t, isNewGame)

	err = table.ReserveSeat(sPosition, secondIdentity)
	assert.Nil(t, err)
	isNewGame, err = table.BuyIn(sChips, sPosition, secondIdentity)
	assert.Nil(t, err)
	assert.True(t, isNewGame)

	return table
}

func table2Player_gameRiver(t *testing.T) *Table {
	return table2PlayerPositions_gameRiver(t, 0, 1)
}

// Returns built table, bb position, sb position
func table2PlayerPositions_gameRiver(t *testing.T, fPosition int, sPosition int) *Table {
	table := table2Players_startedGame(t)
	dealerSbPosition := table.DecidingPosition
	bbPosition := notPlayerPosition(table.DecidingPlayerUnsafe(), fPosition, sPosition)

	err := table.MakeActionDeprecated(dealerSbPosition, Call, 0)
	assert.Nil(t, err)
	err = table.MakeActionDeprecated(bbPosition, Check, 0)
	assert.Nil(t, err)
	// Flop
	err = table.MakeActionDeprecated(bbPosition, Check, 0)
	assert.Nil(t, err)
	err = table.MakeActionDeprecated(dealerSbPosition, Check, 0)
	assert.Nil(t, err)
	// Turn
	err = table.MakeActionDeprecated(bbPosition, Check, 0)
	assert.Nil(t, err)
	err = table.MakeActionDeprecated(dealerSbPosition, Check, 0)
	assert.Nil(t, err)

	return table
}

func defaultTable(t *testing.T) *Table {
	table, err := NewTable(defaultTableParams)
	assert.Nil(t, err)
	return table
}

func tableParams(size int) NewTableParams {
	return NewTableParams{
		Name:            "my table",
		Size:            size,
		BigBlind:        2,
		DecisionTimeout: 0,
		Identity:        firstIdentity,
	}
}

func table3Players_startedGame(t *testing.T) *Table {
	return table3Players_startedGame_chips(t, 0, 1, 2, defaultBuyIn, defaultBuyIn, defaultBuyIn)
}

// If you are using algoMock, then the next deciding player is sPosition
func table3Players_startedGame_chips(t *testing.T, fPosition, sPosition, tPosition int, fChips, sChips, tChips int64) *Table {
	fChips++
	sChips--
	table := defaultTable(t)
	err := table.ReserveSeat(fPosition, firstIdentity)
	assert.Nil(t, err)
	isNewGame, err := table.BuyIn(fChips, fPosition, firstIdentity)
	assert.Nil(t, err)
	assert.False(t, isNewGame)

	err = table.ReserveSeat(sPosition, secondIdentity)
	assert.Nil(t, err)
	isNewGame, err = table.BuyIn(sChips, sPosition, secondIdentity)
	assert.Nil(t, err)
	assert.True(t, isNewGame)

	err = table.ReserveSeat(tPosition, thirdIdentity)
	assert.Nil(t, err)
	newGame, err := table.BuyIn(tChips, tPosition, thirdIdentity)
	assert.Nil(t, err)
	assert.False(t, newGame)

	smallBlindPosition := table.DecidingPosition
	err = table.MakeActionDeprecated(smallBlindPosition, Fold, 0)
	assert.Nil(t, err)

	assert.Nil(t, table.StartNextGame())

	assert.Len(t, table.PlayingPlayersByGameType(), 3)

	return table
}

// Don't use without MockAlgo
// Deciding player is fPosition
func table3Players_startedGameOrder_chips(t *testing.T, fPosition, sPosition, tPosition int, fChips, sChips, tChips int64) *Table {
	table := defaultTable(t)
	err := table.ReserveSeat(sPosition, secondIdentity)
	assert.Nil(t, err)
	isNewGame, err := table.BuyIn(sChips, sPosition, secondIdentity)
	assert.Nil(t, err)
	assert.False(t, isNewGame)

	err = table.ReserveSeat(tPosition, thirdIdentity)
	assert.Nil(t, err)
	isNewGame, err = table.BuyIn(tChips, tPosition, thirdIdentity)
	assert.Nil(t, err)
	assert.True(t, isNewGame)

	assert.Nil(t, table.MakeActionDeprecated(table.DecidingPosition, Fold, 0))
	assert.Nil(t, table.StartNextGame())
	_, _, d := table.BlindsPlayers()
	assert.EqualValues(t, tPosition, d.Position)

	err = table.ReserveSeat(fPosition, firstIdentity)
	assert.Nil(t, err)
	newGame, err := table.BuyIn(fChips, fPosition, firstIdentity)
	assert.Nil(t, err)
	assert.False(t, newGame)

	err = table.MakeActionDeprecated(table.DecidingPosition, Fold, 0)
	assert.Nil(t, err)

	assert.Nil(t, table.StartNextGame())

	assert.Len(t, table.PlayingPlayersByGameType(), 3)

	bb, sb, d := table.BlindsPlayers()
	assert.EqualValues(t, fPosition, d.Position)
	assert.EqualValues(t, sPosition, sb.Position)
	assert.EqualValues(t, tPosition, bb.Position)

	return table
}

func table3Players_flop(t *testing.T) *Table {
	table := table3Players_startedGameOrder_chips(t, 0, 1, 2, defaultBuyIn, defaultBuyIn, defaultBuyIn)
	assert.Nil(t, table.MakeActionDeprecated(table.DecidingPosition, Call, 0))
	assert.Nil(t, table.MakeActionDeprecated(table.DecidingPosition, Call, 0))
	assert.Nil(t, table.MakeActionDeprecated(table.DecidingPosition, Check, 0))
	assert.True(t, table.IsNewRound())
	return table
}

func table4Players_startedGame(t *testing.T) *Table {
	return table4PlayersPositions_startedGame_chips(t, 0, 1, 2, 3, defaultBuyIn, defaultBuyIn, defaultBuyIn, defaultBuyIn)
}

func table4Players_startedGame_chips(t *testing.T, fChips, sChips, tChips, foChips int64) *Table {
	return table4PlayersPositions_startedGame_chips(t, 0, 1, 2, 3, fChips, sChips, tChips, foChips)
}

// Don't use without MockAlgo
func table4PlayersPositions_startedGame_chips(t *testing.T, fPos, sPos, tPos, foPos int, fChips, sChips, tChips, foChips int64) *Table {
	fChips++
	sChips--
	table := defaultTable(t)
	err := table.ReserveSeat(fPos, firstIdentity)
	assert.Nil(t, err)
	isNewGame, err := table.BuyIn(fChips, fPos, firstIdentity)
	assert.Nil(t, err)
	assert.False(t, isNewGame)

	err = table.ReserveSeat(sPos, secondIdentity)
	assert.Nil(t, err)
	isNewGame, err = table.BuyIn(sChips, sPos, secondIdentity)
	assert.Nil(t, err)
	assert.True(t, isNewGame)

	assert.Nil(t, table.MakeActionDeprecated(table.DecidingPosition, Fold, 0))
	assert.Nil(t, table.StartNextGame())
	assert.Nil(t, table.MakeActionDeprecated(table.DecidingPosition, Fold, 0))
	assert.Nil(t, table.StartNextGame())

	err = table.ReserveSeat(tPos, thirdIdentity)
	assert.Nil(t, err)
	_, err = table.BuyIn(tChips, tPos, thirdIdentity)
	assert.Nil(t, err)
	err = table.ReserveSeat(foPos, fourthIdentity)
	assert.Nil(t, err)
	_, err = table.BuyIn(foChips, foPos, fourthIdentity)
	assert.Nil(t, err)

	err = table.MakeActionDeprecated(table.DecidingPosition, Fold, 0)
	assert.Nil(t, err)

	assert.Nil(t, table.StartNextGame())

	assert.Len(t, table.PlayingPlayersByGameType(), 4)

	bb, sb, d := table.BlindsPlayers()
	assert.EqualValues(t, sPos, d.Position)
	assert.EqualValues(t, tPos, sb.Position)
	assert.EqualValues(t, foPos, bb.Position)

	return table
}

func (t *Table) tsit(test *testing.T, pos int, chips int64) {
	assert.Nil(test, t.ReserveSeat(pos, getIden(pos)))
	_, err := t.BuyIn(chips, pos, getIden(pos))
	assert.Nil(test, err)
}

func (t *Table) tsitBack(test *testing.T, pos int) {
	isNewGame, err := t.SitBack(pos, getIden(pos))
	assert.False(test, isNewGame)
	assert.Nil(test, err)
}

func (t *Table) tsetIntent(test *testing.T, pos int, intent Intent) {
	err := t.SetIntent(pos, getIden(pos), intent)
	assert.Nil(test, err)
}

func (t *Table) ttimeout(test *testing.T) {
	assert.Nil(test, t.MakeActionOnTimeout(t.DecidingPosition))
}

func (t *Table) tcheck(test *testing.T) {
	t.taction(test, CheckAction)
}

func (t *Table) tfold(test *testing.T) {
	t.taction(test, FoldAction)
}

func (t *Table) tcall(test *testing.T) {
	t.taction(test, CallAction)
}

func (t *Table) tallIn(test *testing.T) {
	t.taction(test, AllInAction)
}

func (t *Table) traise(test *testing.T, chips int64) {
	t.taction(test, RaiseAction(chips))
}

func (t *Table) tbet(test *testing.T, chips int64) {
	t.taction(test, BetAction(chips))
}

func (t *Table) taction(test *testing.T, action Action) {
	assert.Nil(test, t.MakeAction(t.DecidingPosition, t.DecidingPlayerUnsafe().Identity, action))
}
