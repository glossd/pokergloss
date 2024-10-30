package e2e

import (
	"context"
	"fmt"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/multi"
	"github.com/glossd/pokergloss/table/services/player"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	"github.com/glossd/pokergloss/table/web/client/mq"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestRebalancer_AllMadeAllIn_RemoveTable_P4_TS3(t *testing.T) {
	multiSetUp(t)

	lobby := insertFullLobbyMulti(t, NewLobbyMultiParams{numOfUsers: 4})

	algo, err := domain.NewMockAlgoMultiGame(
		domain.CardsStr("2s", "7d", "Ad", "As"),
		domain.CardsStr("2s", "7d", "Ad", "As"))
	assert.Nil(t, err)
	domain.Algo = algo

	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))
	lobby, err = db.FindLobbyMultiNoCtx(lobby.ID)
	assert.Nil(t, err)

	assert.Len(t, lobby.TableIDs, 2)

	communityCardsMock := domain.NewMockCards("Ac", "Ah", "Kd", "Ks", "Qh")

	domain.Algo = communityCardsMock
	restMakeAction(t, lobby.TableIDs[0].Hex(), domain.AllIn, 0, defaultToken)
	restMakeAction(t, lobby.TableIDs[0].Hex(), domain.AllIn, 1, secondPlayerToken)

	domain.Algo = communityCardsMock
	restMakeAction(t, lobby.TableIDs[1].Hex(), domain.AllIn, 0, thirdPlayerToken)
	restMakeAction(t, lobby.TableIDs[1].Hex(), domain.AllIn, 1, fourthPlayerToken)

	table1 := assertTablePlayersLen(t, lobby.TableIDs[0], 1)
	assert.EqualValues(t, domain.WaitingTable, table1.Status)
	assertTablePlayersLen(t, lobby.TableIDs[1], 1)

	multi.Rebalance(lobby.ID)
	//
	//table1 = assertTablePlayersLen(t, lobby.TableIDs[0], 2)
	//assert.EqualValues(t, domain.PlayingTable, table1.Status)
	//
	//_, err = db.FindTableNoCtx(lobby.TableIDs[1])
	//assert.True(t, errors.Is(err, mongo.ErrNoDocuments))
}

func TestRebalancer_WaitToRemoveForGameEnd_P4_TS3(t *testing.T) {
	multiSetUp(t)

	lobby := insertFullLobbyMulti(t, NewLobbyMultiParams{numOfUsers: 4})

	algo, err := domain.NewMockAlgoMultiGame(
		domain.CardsStr("2s", "7d", "Ad", "As"),
		domain.CardsStr("2s", "7d", "Ad", "As"))
	assert.Nil(t, err)
	domain.Algo = algo

	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))
	lobby = findLobby(t, lobby.ID, 2)

	assertMultiGameStartSimple(t)
	assertMultiLobby(t, lobby)

	assert.Len(t, lobby.TableIDs, 2)

	communityCardsMock := domain.NewMockCards("Ac", "Ah", "Kd", "Ks", "Qh")

	domain.Algo = communityCardsMock
	restMakeAction(t, lobby.TableIDs[0].Hex(), domain.AllIn, 0, defaultToken)
	assertSimpleAction(t, 0, 1, domain.AllIn)
	restMakeAction(t, lobby.TableIDs[0].Hex(), domain.AllIn, 1, secondPlayerToken)
	assertActionGameEndWithShowdown(t, 1, domain.AllIn, 0, 1, 1)
	assertMessage(t, 2, func(as []*Asserter) {
		as[0].assertPlayerLeft(0)
		as[1].assertReset(3)
	})

	multi.Rebalance(lobby.ID)

	assertMultiPlayersUpdate(t, lobby.TableIDs[0].Hex(), 0)

	assertMultiPlayerMove(t, []string{secondIdentity.UserId}, lobby.TableIDs[1].Hex())

	assertTableDontExist(t, lobby.TableIDs[0])
	table2Id := lobby.TableIDs[1].Hex()
	restMakeAction(t, table2Id, domain.Fold, 0, thirdPlayerToken)

	table2 := findTable(t, lobby.TableIDs[1])
	assert.EqualValues(t, domain.PlayingTable, table2.Status)
	assert.EqualValues(t, 1, table2.DealerPosition())
	assert.EqualValues(t, 2, table2.SmallBlindPosition())
	assert.EqualValues(t, 0, table2.BigBlindPosition())
	assert.True(t, table2.MultiAttrs.IsLast)

	restMakeAction(t, table2Id, domain.Call, 1, fourthPlayerToken)
}

func TestRebalancer_MoveTwoPlayers_And_DeleteTable(t *testing.T) {
	multiSetUp(t)

	lobby := insertFullLobbyMulti(t, NewLobbyMultiParams{numOfUsers: 6, tableSize: 5})
	domain.Algo = domain.NewMockCards("As", "Ad", "Ks", "8d", "2s", "7h") // first table hole cards
	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))
	lobby = findLobby(t, lobby.ID, 2)

	assertMultiGameStartSimple(t)
	assertMultiLobby(t, lobby)

	domain.Algo = domain.NewMockCards("Ah", "Ac", "Kh", "Kc", "9s") // community cards of first table
	hex1 := lobby.TableIDs[0].Hex()
	restMakeAction(t, hex1, domain.AllIn, 0, getToken(0))
	assertSimpleAction(t, 0, 1, domain.AllIn)
	restMakeAction(t, hex1, domain.Fold, 1, getToken(1))
	assertSimpleAction(t, 1, 2, domain.Fold)
	restMakeAction(t, hex1, domain.AllIn, 2, getToken(2))
	readMessage() // betting action winners
	readMessage() // todo assert player left and start hand
	assertMultiPlayersUpdate(t, hex1, 2)

	multi.Rebalance(lobby.ID)

	restMakeAction(t, hex1, domain.Fold, 1, getToken(1))
	msg := readMessage() // end game, showdown, winners

	hex2 := lobby.TableIDs[1].Hex()
	msg = readMessage()
	NewAsserter(t, msg.UserEvents.UserEvents[defaultIdentity.UserId].Events[0]).assertMultiPlayerMove(hex2)
	NewAsserter(t, msg.UserEvents.UserEvents[secondIdentity.UserId].Events[0]).assertMultiPlayerMove(hex2)
	playerLeftEvents := msg.UserEvents.AfterEvents
	assert.EqualValues(t, 2, len(playerLeftEvents))
	assert.EqualValues(t, 1, gjson.Get(playerLeftEvents[0].Payload, "leftPlayer.position").Int())
	assert.EqualValues(t, 0, gjson.Get(playerLeftEvents[1].Payload, "leftPlayer.position").Int())

	table2 := findTable(t, lobby.TableIDs[1])

	assert.EqualValues(t, 5, len(table2.AllPlayers()))

	multi.Rebalance(lobby.ID) // should delete table1

	assertTableDontExist(t, lobby.TableIDs[0])
	table2 = findTable(t, lobby.TableIDs[1])
	assert.True(t, table2.IsLast)

	restMakeAction(t, hex2, domain.Fold, 0, fourthPlayerToken)
	restMakeAction(t, hex2, domain.Fold, 1, fifthPlayerToken)
	// new game
	table2 = findTable(t, lobby.TableIDs[1])
	assert.True(t, table2.IsLast)
}

func findLobby(t *testing.T, lobbyID primitive.ObjectID, tablesNum int) *domain.LobbyMulti {
	lobby, err := db.FindLobbyMultiNoCtx(lobbyID)
	assert.Nil(t, err)
	assert.Len(t, lobby.TableIDs, tablesNum)
	return lobby
}

func Test_6P_5P_5PlayersGoAllInAndOnly2Left(t *testing.T) {
	multiSetUp(t)

	lobby := insertFullLobbyMulti(t, NewLobbyMultiParams{numOfUsers: 11, tableSize: 6})
	domain.Algo = domain.NewMockCardsSkip(12,
		"As", "Ks",
		"Ad", "Kd",
		"3s", "7d",
		"4h", "8d",
		"2s", "4c",
	) // second table hole cards
	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))
	lobby, err := db.FindLobbyMultiNoCtx(lobby.ID)
	assert.Nil(t, err)
	assert.Len(t, lobby.TableIDs, 2)

	domain.Algo = domain.NewMockCards("Ah", "Ac", "Kh", "Kc", "9s") // community cards of second table
	hex2 := lobby.TableIDs[1].Hex()
	restMakeAction(t, hex2, domain.AllIn, 3, getToken2Table(3))
	restMakeAction(t, hex2, domain.AllIn, 4, getToken2Table(4))
	restMakeAction(t, hex2, domain.AllIn, 0, getToken2Table(0))
	restMakeAction(t, hex2, domain.AllIn, 1, getToken2Table(1))
	restMakeAction(t, hex2, domain.AllIn, 2, getToken2Table(2))

	assertTablePlayersLen(t, lobby.TableIDs[1], 2)
	multi.Rebalance(lobby.ID)

	hex1 := lobby.TableIDs[0].Hex()
	restMakeAction(t, hex1, domain.Fold, 3, getToken(3))
	restMakeAction(t, hex1, domain.Fold, 4, getToken(4))
	restMakeAction(t, hex1, domain.Fold, 5, eleventhPlayerToken)
	restMakeAction(t, hex1, domain.Fold, 0, getToken(0))
	restMakeAction(t, hex1, domain.Fold, 1, getToken(1))

	assertTablePlayersLen(t, lobby.TableIDs[0], 4)
	assertTablePlayersLen(t, lobby.TableIDs[1], 4)
}

func TestRebalancer_CountDown_RemoveTable(t *testing.T) {
	multiCountDownSetUp(t)

	lobby := insertFullLobbyMulti(t, NewLobbyMultiParams{numOfUsers: 4})

	algo, err := domain.NewMockAlgoMultiGame(domain.CardsStr("2s", "7d", "Ad", "As"))
	assert.Nil(t, err)
	domain.Algo = algo

	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))
	lobby, err = db.FindLobbyMultiNoCtx(lobby.ID)
	assert.Nil(t, err)

	assert.Len(t, lobby.TableIDs, 2)

	communityCardsMock := domain.NewMockCards("Ac", "Ah", "Kd", "Ks", "Qh")

	domain.Algo = communityCardsMock
	hex0 := lobby.TableIDs[0].Hex()
	restMakeAction(t, hex0, domain.AllIn, 0, defaultToken)
	restMakeAction(t, hex0, domain.AllIn, 1, secondPlayerToken)

	time.Sleep(10 * time.Millisecond)

	assertTableDontExist(t, lobby.TableIDs[0])
	assertTablePlayersLen(t, lobby.TableIDs[1], 3)
}

func multiCountDownSetUp(t *testing.T) {
	multiSetUp(t)
	conf.Props.RebalancerPeriod = 0
	mq.IsMultiPlayersMovedEnabledTest = true
	t.Cleanup(func() {
		mq.IsMultiPlayersMovedEnabledTest = false
	})
}

func TestRebalancer_CountDown_MoveAllPlayers(t *testing.T) {
	multiCountDownSetUp(t)

	lobby := insertFullLobbyMulti(t, NewLobbyMultiParams{numOfUsers: 5, tableSize: 4})

	domain.Algo = domain.NewMockCardsSkip(6, "2s", "7d", "Ad", "As")

	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))
	lobby, err := db.FindLobbyMultiNoCtx(lobby.ID)
	assert.Nil(t, err)

	assert.Len(t, lobby.TableIDs, 2)

	communityCardsMock := domain.NewMockCards("Ac", "Ah", "Kd", "Ks", "Qh")

	domain.Algo = communityCardsMock
	hex1 := lobby.TableIDs[1].Hex()
	restMakeAction(t, hex1, domain.AllIn, 0, thirdPlayerToken)
	restMakeAction(t, hex1, domain.AllIn, 1, fourthPlayerToken)

	time.Sleep(20 * time.Millisecond)

	assertTableDontExist(t, lobby.TableIDs[1])
	assertTablePlayersLen(t, lobby.TableIDs[0], 4)
}

func TestRebalancer_CountDown_Disproportion(t *testing.T) {
	multiCountDownSetUp(t)

	lobby := insertFullLobbyMulti(t, NewLobbyMultiParams{numOfUsers: 7, tableSize: 4})

	domain.Algo = domain.NewMockCardsSkip(8, "2s", "7d", "Ad", "As")

	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))
	lobby, err := db.FindLobbyMultiNoCtx(lobby.ID)
	assert.Nil(t, err)

	assert.Len(t, lobby.TableIDs, 2)

	communityCardsMock := domain.NewMockCards("Ac", "Ah", "Kd", "Ks", "Qh")

	domain.Algo = communityCardsMock
	hex1 := lobby.TableIDs[1].Hex()
	restMakeAction(t, hex1, domain.AllIn, 0, fourthPlayerToken)
	restMakeAction(t, hex1, domain.AllIn, 1, fifthPlayerToken)
	restMakeAction(t, hex1, domain.Fold, 2, sixthPlayerToken)

	assertTablePlayersLen(t, lobby.TableIDs[1], 2)
	assertTablePlayersLen(t, lobby.TableIDs[0], 4)

	hex0 := lobby.TableIDs[0].Hex()
	restMakeAction(t, hex0, domain.Fold, 3, seventhPlayerToken)
	restMakeAction(t, hex0, domain.Fold, 0, defaultToken)
	restMakeAction(t, hex0, domain.Fold, 1, secondPlayerToken)

	assertTablePlayersLen(t, lobby.TableIDs[1], 3)
	assertTablePlayersLen(t, lobby.TableIDs[0], 3)
}

func TestRebalancer_MovePlayers_RaceCondition(t *testing.T) {
	multiCountDownSetUp(t)
	mq.IsMultiPlayersMovedEnabledTest = false
	conf.Props.Table.GameEndMinTimeout = -1

	lobby := insertFullLobbyMulti(t, NewLobbyMultiParams{numOfUsers: 7, tableSize: 6})

	domain.Algo = domain.NewMockCardsSkip(8, "2s", "7d", "Ad", "As")

	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))
	lobby, err := db.FindLobbyMultiNoCtx(lobby.ID)
	assert.Nil(t, err)

	assert.Len(t, lobby.TableIDs, 2)

	communityCardsMock := domain.NewMockCards("Ac", "Ah", "Kd", "Ks", "Qh")

	domain.Algo = communityCardsMock
	hex1 := lobby.TableIDs[1].Hex()
	restMakeAction(t, hex1, domain.AllIn, 0, fourthPlayerToken)
	restMakeAction(t, hex1, domain.AllIn, 1, fifthPlayerToken)
	restMakeAction(t, hex1, domain.Fold, 2, sixthPlayerToken)
	doStartGame(t, lobby.TableIDs[1])

	result, err := multi.Rebalance(lobby.ID)
	assert.Nil(t, err)
	assert.EqualValues(t, multi.MoveAllPlayers, result.Status)
	restMakeAction(t, hex1, domain.Fold, 1, fifthPlayerToken)

	table0 := findTable(t, lobby.TableIDs[0])
	// race condition
	doStartGame(t, lobby.TableIDs[1])

	params, err := player.NewPositionParams(context.Background(), lobby.TableIDs[0].Hex(), 3, seventhIdentity)
	assert.Nil(t, err)
	chipsParams, err := player.ToChipsParams(params, 0)
	assert.Nil(t, err)
	fmt.Println(table0.DecidingPosition)
	err = player.DoBettingActionOnTable(table0, chipsParams, domain.Fold)
	assert.Nil(t, err)

	assertTablePlayersLen(t, lobby.TableIDs[1], 0)
	assertTablePlayersLen(t, lobby.TableIDs[0], 6)
}

func doStartGame(t *testing.T, id primitive.ObjectID) {
	table := findTable(t, id)
	actionhandler.DoStartGameNoCtx(timeout.Key{TableID: table.ID, Version: table.GameFlowVersion})
}
