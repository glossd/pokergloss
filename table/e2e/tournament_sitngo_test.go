package e2e

import (
	"encoding/json"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/model"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"testing"
	"time"
)

func TestSitAndGoLobby_Post(t *testing.T) {
	t.Cleanup(cleanUp)

	body := postSingAndGoLobby(t)

	var postRes model.LobbySitAndGo
	err := json.Unmarshal([]byte(body), &postRes)
	assert.Nil(t, err)
}

func TestSitAndGoLobby_List(t *testing.T) {
	body := listSitAndGoLobbies(t)
	assert.EqualValues(t, "[]", body)

	for i := 0; i < 5; i++ {
		postSingAndGoLobby(t)
	}

	body = listSitAndGoLobbies(t)
	var res []model.LobbySitAndGo
	err := json.Unmarshal([]byte(body), &res)
	assert.Nil(t, err)
	assert.Len(t, res, 5)
}

func listSitAndGoLobbies(t *testing.T) string {
	rr := testRouter.Request(t, http.MethodGet, "/sit-n-go/lobbies", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	return body
}

func postSingAndGoLobby(t *testing.T) string {
	rr := testRouter.Request(t, http.MethodPost, "/sit-n-go/lobbies", postSitAndGoLobbyBody(), nil)

	body := rr.Body.String()
	assert.Equal(t, http.StatusCreated, rr.Code, body)
	return body
}

func postSitAndGoLobbyBody() *string {
	b := `{"tableParams": {"name":"my table", "size":2, "bigBlind":2}, "placesPaid": 1, "buyIn": 2500, "levelIncreaseTimeMinutes":1, "startingChips": 250}`
	return &b
}

func TestSitAndGoLobby_Register(t *testing.T) {
	t.Cleanup(cleanUp)

	lobby := insertLobbySitAndGo(t)

	putBody := `{"position":1}`
	rr := testRouter.Request(t, http.MethodPut, "/sit-n-go/lobbies/"+lobby.ID.Hex()+"/register", &putBody, authHeaders(secondPlayerToken))
	body := rr.Body.String()
	assert.EqualValues(t, http.StatusOK, rr.Code, body)
}

func TestSitAndGoLobby_Unregister(t *testing.T) {
	t.Cleanup(cleanUp)

	lobby := insertLobbySitAndGo(t)

	putBody := `{"position":0}`
	rr := testRouter.Request(t, http.MethodPut, "/sit-n-go/lobbies/"+lobby.ID.Hex()+"/register", &putBody, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())

	rr = testRouter.Request(t, http.MethodPut, "/sit-n-go/lobbies/"+lobby.ID.Hex()+"/unregister", &putBody, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())
}

func TestSitAndGoLobby_FullGame(t *testing.T) {
	defaultDuration := domain.MinLevelIncreaseMinDuration
	t.Cleanup(func() {
		cleanUp()
		domain.MinLevelIncreaseMinDuration = defaultDuration
		conf.Props.Table.GameEndMinTimeout = defaultGameEndTimeout
	})
	domain.MinLevelIncreaseMinDuration = 0
	conf.Props.Table.GameEndMinTimeout = 0

	lobby := insertFullLobbySitAndGo(t, func() *domain.LobbySitAndGo {
		return initLobbySitngoLevelIncrease(t, 0)
	})

	table := lobby.GetTable()
	restMakeAction(t, table.ID.Hex(), domain.Fold, 0)

	table, err := db.FindTableNoCtx(table.ID)
	assert.Nil(t, err)
	assert.EqualValues(t, 4, table.BigBlind)
}

func TestSitAndGoLobby_2P_AllSittingOut(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.Table.GameEndMinTimeout = 0
	conf.Props.MinDecisionTimeout = -1

	domain.Algo = Algo_2P_MockZeroPositionLoses(t)

	lobby := insertFullLobbySitAndGo(t, func() *domain.LobbySitAndGo {
		return initLobbySitngoTableParams(t, tableParams(2, -1))
	})

	table := lobby.GetTable()
	actionhandler.DoDecisionTimeoutNoCtx(timeout.Key{TableID: table.ID, Position: 0, Version: table.GameFlowVersion})
	assertMessage(t, 4, func(as []*Asserter) {
		as[0].assertTimeToDecideTimeout(0)
		as[1].assertStackOverFlowPlayer()
		as[2].assertShowdown(1, true)
		as[3].assertWinners(1)
	})

	assertStartHandTableSize(t, 0, 1, 1, 2)

	table, err := db.FindTableNoCtx(table.ID)
	assert.Nil(t, err)
	actionhandler.DoDecisionTimeoutNoCtx(timeout.Key{TableID: table.ID, Position: 1, Version: table.GameFlowVersion})
	assertMessage(t, 3, func(as []*Asserter) {
		as[0].assertTimeToDecideTimeout(1)
		as[1].assertStackOverFlowPlayer()
		as[2].assertWinners(1)
	})

	assertMessage(t, 3, func(as []*Asserter) {
		as[0].assertPlayerLeftPrize(0, 2, 0)
		as[1].assertPlayerLeftPrize(1, 1, defaultBuyIn*2)
		as[2].assertReset(2)
	})

	_, err = db.FindTableNoCtx(table.ID)
	assert.EqualValues(t, err, mongo.ErrNoDocuments)
}

func TestSitAndGoLobby_FirstHand_HoleCardsNotFaceDown(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.Table.GameEndMinTimeout = 0
	conf.Props.MinDecisionTimeout = -1

	lobby := insertLobbySitAndGo(t)
	restSitngoRegister(t, lobby.ID.Hex(), 0, getToken(0))
	assertSitngoRegister(t, 0)

	restSitngoRegister(t, lobby.ID.Hex(), 1, getToken(1))
	assertSitngoRegister(t, 1)
	assertSitngoGameStartSimple(t)

	lobby, err := db.FindSitAndGoLobbyNoCtx(lobby.ID)
	assert.Nil(t, err)

	table := restGetFullTable(t, lobby.TableID.Hex(), getToken(0))
	v := table.Seats[0].Player.Cards
	holeCards := *v
	assert.Len(t, holeCards, 2)
	assert.NotEmpty(t, holeCards[0])
	assert.NotEqualValues(t, "Xx", holeCards[0])
}

func TestSitAndGoLobby_DecisionTimeoutOfFirstPlayer(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.MinDecisionTimeout = -1
	domain.Algo = &domain.MockAlgo{}

	lobby := initLobbySitngoTableParams(t, domain.NewTableParams{Name: "sitngo", Size: 2, DecisionTimeout: 1, BigBlind: 2, Identity: defaultIdentity})
	err := db.InsertSitAndGoLobbyNoCtx(lobby)
	assert.Nil(t, err)

	restSitngoRegister(t, lobby.ID.Hex(), 0, getToken(0))
	assertSitngoRegister(t, 0)

	restSitngoRegister(t, lobby.ID.Hex(), 1, getToken(1))
	assertSitngoRegister(t, 1)

	assertTimeToDecideTimeoutAndStackOverFlowAndWinners(t, 0, 1) // decisionTimeout gets first in TestMQ
}

func insertLobbySitAndGo(t *testing.T) *domain.LobbySitAndGo {
	lobby := initLobbySitngo(t)

	err := db.InsertSitAndGoLobbyNoCtx(lobby)
	assert.Nil(t, err)
	return lobby
}

func initLobbySitngo(t *testing.T) *domain.LobbySitAndGo {
	return initLobbySitngoLevelIncrease(t, time.Minute)
}

func initLobbySitngoLevelIncrease(t *testing.T, levelIncreaseTime time.Duration) *domain.LobbySitAndGo {
	params := domain.NewLobbySitAndGoParams{
		NewTableParams:    tableParams(2, -1),
		PlacesPaid:        1,
		BuyIn:             5000,
		LevelIncreaseTime: levelIncreaseTime,
	}

	lobby, err := domain.NewLobbySitAndGo(params)
	assert.Nil(t, err)
	return lobby
}

func initLobbySitngoTableParams(t *testing.T, tableParams domain.NewTableParams) *domain.LobbySitAndGo {
	params := domain.NewLobbySitAndGoParams{
		NewTableParams:    tableParams,
		PlacesPaid:        1,
		BuyIn:             defaultBuyIn,
		LevelIncreaseTime: time.Minute,
	}

	lobby, err := domain.NewLobbySitAndGo(params)
	assert.Nil(t, err)
	return lobby
}

func insertFullLobbySitAndGo(t *testing.T, initLobby func() *domain.LobbySitAndGo) *domain.LobbySitAndGo {
	lobby := initLobby()
	assert.Nil(t, lobby.Register(defaultIdentity, 0))
	assert.Nil(t, lobby.Register(secondIdentity, 1))
	assert.Nil(t, db.InsertSitAndGoLobbyNoCtx(lobby))
	insertTable(t, lobby.GetTable())
	return lobby
}
