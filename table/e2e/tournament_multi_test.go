package e2e

import (
	"context"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/multi"
	"github.com/glossd/pokergloss/table/services/paging"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

const defaultTournamentBuyIn = 250
const defaultTournamentBigBling = 2

var defaultMultiIncrease = 5 * time.Minute

func TestCreateDailyFreerolls(t *testing.T) {
	t.Cleanup(cleanUp)

	multi.CreateDailyMultiTournaments(time.Now().AddDate(0, 0, 1), multi.NothingEnrich)
	lobbies, err := db.FindAllMultiLobbiesNoCtx()
	assert.Nil(t, err)
	assert.Positive(t, len(lobbies))
}

func TestStartLobbyMulti_NotEnoughPlayers(t *testing.T) {
	t.Cleanup(cleanUp)

	lobby := defaultLobbyMulti()
	assert.Nil(t, db.InsertOneLobbyMulti(lobby))

	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	lobbies, err := db.FindMultiLobbies(ctx, paging.DefaultParams)
	assert.Nil(t, err)
	assert.Zero(t, len(lobbies))
}

func TestStartLobbyMulti(t *testing.T) {
	prevPropsSetup(t)

	lobby := insertFullLobbyMulti(t, NewLobbyMultiParams{numOfUsers: 2, decisionTimeout: 1})

	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))
	lobbies := findLobbies(t, lobby.StartAt)
	lobby = lobbies[0]
	assertMultiGameStart(t, lobby.TableIDs[0].Hex())
	assertMultiLobby(t, lobby)

	assertTimeToDecideTimeoutAndStackOverFlowAndWinners(t, 0, 1)
}

func TestMultiLobby_IncreaseBlinds(t *testing.T) {
	multiSetUp(t)
	defaultMultiIncrease = 0

	lobby := insertFullLobbyMultiDeprecated(t, 2, 2)
	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))
	lobbies := findLobbies(t, lobby.StartAt)
	assert.Len(t, lobbies, 1)
	lobby = lobbies[0]
	assert.Len(t, lobby.TableIDs, 1)

	table := findTable(t, lobby.TableIDs[0])
	assert.EqualValues(t, 2, table.BigBlind)
	assert.EqualValues(t, 2, table.NextSmallBlind)

	restMakeAction(t, table.ID.Hex(), domain.Fold, 0)

	table = findTable(t, table.ID)
	assert.EqualValues(t, 4, table.BigBlind)
	assert.EqualValues(t, 3, table.NextSmallBlind)
}

func TestMultiLobby_FinishTournament(t *testing.T) {
	multiSetUp(t)

	lobby := insertFullLobbyMulti(t, NewLobbyMultiParams{})

	domain.Algo = domain.NewMockCards("2s", "7d", "Ad", "As")
	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))
	lobbies := findLobbies(t, lobby.StartAt)
	assert.Len(t, lobbies, 1)
	lobby = lobbies[0]

	domain.Algo = domain.NewMockCards("Ac", "Ah", "Kd", "Ks", "Qh")
	restMakeAction(t, lobby.TableIDs[0].Hex(), domain.AllIn, 0, defaultToken)
	restMakeAction(t, lobby.TableIDs[0].Hex(), domain.AllIn, 1, secondPlayerToken)

	assertTableDontExist(t, lobby.TableIDs[0])
}

func TestMultiLobby_ShouldDeleteTournament_WhenNotEnoughRegistered(t *testing.T) {
	multiSetUp(t)

	lobby := lobbyMulti(NewLobbyMultiParams{})
	assert.Nil(t, db.InsertOneLobbyMulti(lobby))

	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))
	lobbies := findLobbies(t, lobby.StartAt)
	assert.Len(t, lobbies, 0)
}

func TestMultiLobby_ComputeWinners(t *testing.T) {
	multiSetUp(t)

	lobby := insertFullLobbyMultiDeprecated(t, 3, 2)

	domain.Algo = domain.NewMockCards("2s", "7d", "Ad", "As")
	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))
	lobbies := findLobbies(t, lobby.StartAt)
	assert.Len(t, lobbies, 1)
	lobby = lobbies[0]
	assertMultiGameStart(t, lobby.TableIDs[0].Hex())
	assertMultiLobby(t, lobby)

	domain.Algo = domain.NewMockCards("Ac", "Ah", "Kd", "Ks", "Qh")
	// todo fix multi players on start
	//assertMultiPlayersUpdate(t, lobby.TableIDs[0].Hex(), 2)

	restMakeAction(t, lobby.TableIDs[0].Hex(), domain.AllIn, 0, defaultToken)
	assertSimpleAction(t, 0, 1, domain.AllIn)

	restMakeAction(t, lobby.TableIDs[0].Hex(), domain.AllIn, 1, secondPlayerToken)
	readMessage()

	assertMessage(t, 3, func(as []*Asserter) {
		as[0].assertPlayerLeftPrize(0, 2, 0)
		as[1].assertPlayerLeftPrize(1, 1, 2*defaultTournamentBuyIn)
		//as[2].assertReset(0) // todo somehow seats are nil
	})
}

// 03.04.2021 bug that killed 10 min of tournament
func TestMultiLobby_OneOfTablesAllSittingOut(t *testing.T) {
	multiSetUp(t)

	domain.Algo = domain.NewMockFull(0)

	lobby := insertFullLobbyMulti(t, NewLobbyMultiParams{tableSize: 6, numOfUsers: 10})
	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))
	lobby, err := db.FindLobbyMultiNoCtx(lobby.ID)
	assert.Nil(t, err)
	assert.Len(t, lobby.TableIDs, 2)

	oid0 := lobby.TableIDs[0]
	assert.False(t, actionhandler.DoDecisionTimeoutNoCtx(timeout.Key{TableID: oid0, Position: 3, Version: 0}))
	assert.False(t, actionhandler.DoDecisionTimeoutNoCtx(timeout.Key{TableID: oid0, Position: 4, Version: 1}))
	assert.False(t, actionhandler.DoDecisionTimeoutNoCtx(timeout.Key{TableID: oid0, Position: 0, Version: 2}))
	restMakeBetAction(t, oid0.Hex(), domain.Raise, 5, 1, getToken(1))
	assert.False(t, actionhandler.DoDecisionTimeoutNoCtx(timeout.Key{TableID: oid0, Position: 2, Version: 4}))
	assert.False(t, actionhandler.DoDecisionTimeoutNoCtx(timeout.Key{TableID: oid0, Position: 0, Version: 5}))
	conf.Props.Table.GameEndMinTimeout = -1
	restMakeAction(t, oid0.Hex(), domain.Fold, 1, secondPlayerToken)
}

func TestMultiLobby_StartMultipleLobbies(t *testing.T) {
	multiSetUp(t)

	now := timeutil.NowAdd(time.Second)
	insertFullLobbyMulti(t, NewLobbyMultiParams{startAt: now})
	insertFullLobbyMulti(t, NewLobbyMultiParams{startAt: now})
	assert.Nil(t, multi.LaunchMultiLobbies(now))
	lobbies := findLobbies(t, now)
	for _, lobby := range lobbies {
		assert.EqualValues(t, domain.LobbyStarted, lobby.Status)
	}
}

func multiSetUp(t *testing.T) {
	increase := defaultMultiIncrease
	prevPropsSetup(t)
	t.Cleanup(func() {
		defaultMultiIncrease = increase
	})
	conf.Props.Multi.RebalancerPeriod = 0
	conf.Props.Table.GameEndMinTimeout = 0
}

func prevPropsSetup(t *testing.T) {
	var prevProps = conf.Props
	t.Cleanup(func() {
		cleanUp()
		conf.Props = prevProps
	})
}

func findTable(t *testing.T, id primitive.ObjectID) *domain.Table {
	table, err := db.FindTableNoCtx(id)
	assert.Nil(t, err)
	return table
}

func findLobbies(t *testing.T, startAt int64) []*domain.LobbyMulti {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	lobbies, err := db.FindMultiLobbiesByStartAt(ctx, startAt)
	assert.Nil(t, err)
	return lobbies
}

func assertTablePlayersLen(t *testing.T, tableID primitive.ObjectID, length int) *domain.Table {
	table, err := db.FindTableNoCtx(tableID)
	assert.Nil(t, err)
	assert.Len(t, table.AllPlayers(), length)
	return table
}

// Deprecated
func insertFullLobbyMultiDeprecated(t *testing.T, tableSize, numOfUsers int) *domain.LobbyMulti {
	lobby := lobbyMulti(NewLobbyMultiParams{tableSize: tableSize, numOfUsers: numOfUsers})
	for i := 0; i < numOfUsers; i++ {
		assert.Nil(t, lobby.Register(getIden(i)))
	}
	assert.Nil(t, db.InsertOneLobbyMulti(lobby))
	return lobby
}

type NewLobbyMultiParams struct {
	tableSize       int // defaults to 3
	startAt         int64
	numOfUsers      int // defaults to 2
	decisionTimeout time.Duration
}

func insertFullLobbyMulti(t *testing.T, params NewLobbyMultiParams) *domain.LobbyMulti {
	lobby := lobbyMulti(params)
	numOfUsers := 2
	if params.numOfUsers > 0 {
		numOfUsers = params.numOfUsers
	}
	for i := 0; i < numOfUsers; i++ {
		assert.Nil(t, lobby.Register(getIden(i)))
	}
	assert.Nil(t, db.InsertOneLobbyMulti(lobby))
	return lobby
}

func lobbyMulti(input NewLobbyMultiParams) *domain.LobbyMulti {
	startAt := timeutil.NowAdd(time.Minute)
	if input.startAt > 0 {
		startAt = input.startAt
	}
	tableSize := 3
	if input.tableSize > 0 {
		tableSize = input.tableSize
	}
	decisionTimeout := time.Duration(-1)
	if input.decisionTimeout >= 0 {
		decisionTimeout = input.decisionTimeout
	}

	params := domain.NewLobbyMultiParams{
		StartAt:           startAt,
		BuyIn:             defaultTournamentBuyIn,
		BigBlind:          defaultTournamentBigBling,
		Name:              "my multi",
		TableSize:         tableSize,
		LevelIncreaseTime: defaultMultiIncrease,
		DecisionTimeout:   decisionTimeout,
	}
	return domain.NewLobbyMulti(params)
}

func defaultLobbyMulti() *domain.LobbyMulti {
	return lobbyMulti(NewLobbyMultiParams{tableSize: 6})
}

func assertMultiPlayersUpdate(t *testing.T, tableID string, playersCount int) {
	assertMessage(t, 1, func(as []*Asserter) {
		as[0].assertMultiPlayersUpdate(tableID, playersCount)
	})
}

func assertMultiPlayerMove(t *testing.T, userIds []string, toTableID string) {
	msg := readMessage()

	assert.EqualValues(t, len(userIds), len(msg.UserEvents.UserEvents))
	for _, userId := range userIds {
		assert.EqualValues(t, 1, len(msg.UserEvents.UserEvents[userId].Events))
		NewAsserter(t, msg.UserEvents.UserEvents[userId].Events[0]).assertMultiPlayerMove(toTableID)
	}
}

func assertMultiLobby(t *testing.T, lobby *domain.LobbyMulti) {
	assertMessage(t, 1, func(as []*Asserter) {
		as[0].assertMultiLobby(lobby)
	})
}
