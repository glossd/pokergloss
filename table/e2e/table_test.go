package e2e

import (
	"encoding/json"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/model"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTable(t *testing.T) {
	t.Cleanup(cleanUp)

	rr := postTable(t)
	b := rr.Body.String()
	id := gjson.Get(b, "id").String()

	tableID, err := primitive.ObjectIDFromHex(id)
	assert.Nil(t, err)
	table, err := db.FindTableNoCtx(tableID)
	assert.Nil(t, err)
	assert.NotEmpty(t, table.CreatedAt)
	assert.EqualValues(t, defaultIdentity.UserId, table.CreateBy)
}

func postTable(t *testing.T) *httptest.ResponseRecorder {
	postBody := `{"name":"my table", "size":9, "bigBlind":2}`
	rr := testRouter.Request(t, http.MethodPost, "/tables", &postBody, nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	return rr
}

func TestGetTables_ShouldReturnEmptyArray(t *testing.T) {
	t.Cleanup(cleanUp)

	rr := testRouter.Request(t, http.MethodGet, "/tables", nil, nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	b := rr.Body.String()
	assert.Equal(t, "[]", b)
}

func TestGetTables_ShouldReturnSameNumberOfTables(t *testing.T) {
	t.Cleanup(cleanUp)

	tableAmount := 4
	for i := 0; i < tableAmount; i++ {
		InsertTable(t)
	}

	rr := testRouter.Request(t, http.MethodGet, "/tables", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	jsonBody := rr.Body.String()
	assert.EqualValues(t, tableAmount, gjson.Get(jsonBody, "#").Int())
}

func TestGetTables_ShouldSkipTables(t *testing.T) {
	t.Cleanup(cleanUp)

	tableAmount := 4
	for i := 0; i < tableAmount; i++ {
		InsertTable(t)
	}

	rr := testRouter.Request(t, http.MethodGet, "/tables?skip=2&limit=4", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())
	assert.EqualValues(t, 2, gjson.Get(rr.Body.String(), "#").Int())

	rr = testRouter.Request(t, http.MethodGet, "/tables?skip=5&limit=4", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 0, gjson.Get(rr.Body.String(), "#").Int())

	rr = testRouter.Request(t, http.MethodGet, "/tables?limit=2", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 2, gjson.Get(rr.Body.String(), "#").Int())
}

func TestGetTables_SkipEmpty(t *testing.T) {
	t.Cleanup(cleanUp)
	tableCount := 4
	tables := make([]*domain.Table, tableCount)
	for i := 0; i < tableCount; i++ {
		tables[i] = InsertTable(t)
	}

	rr := testRouter.Request(t, http.MethodGet, "/tables?skipEmpty=true", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 0, gjson.Get(rr.Body.String(), "#").Int())

	rr = testRouter.Request(t, http.MethodGet, "/tables?skipEmpty=false", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 4, gjson.Get(rr.Body.String(), "#").Int())

	restReserveSeat(t, tables[0].ID.Hex(), 0)

	rr = testRouter.Request(t, http.MethodGet, "/tables?skipEmpty=true", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 1, gjson.Get(rr.Body.String(), "#").Int())
}

func TestGetTables_SkipFull(t *testing.T) {
	t.Cleanup(cleanUp)
	tableCount := 4
	tables := make([]*domain.Table, tableCount)
	for i := 0; i < tableCount-1; i++ {
		tables[i] = InsertTable(t)
	}
	params := domain.NewTableParams{
		Name:            "my table",
		Size:            2,
		BigBlind:        2,
		DecisionTimeout: 0,
		Identity:        defaultIdentity,
	}
	table, err := domain.NewTable(params)
	assert.Nil(t, err)
	insertTable(t, table)
	tables[tableCount-1] = table

	rr := testRouter.Request(t, http.MethodGet, "/tables?skipFull=true", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 4, gjson.Get(rr.Body.String(), "#").Int())

	rr = testRouter.Request(t, http.MethodGet, "/tables?skipFull=false", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 4, gjson.Get(rr.Body.String(), "#").Int())

	restReserveSeat(t, table.ID.Hex(), 0)
	restReserveSeat(t, table.ID.Hex(), 1, secondPlayerToken)

	rr = testRouter.Request(t, http.MethodGet, "/tables?skipFull=true", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 3, gjson.Get(rr.Body.String(), "#").Int())
}

// fields of domain.Card were private, db couldn't save them
func TestSaveCommunityCards(t *testing.T) {
	t.Cleanup(cleanUp)
	domain.Algo = &domain.MockAlgo{}

	sbPosition := 2
	bbPosition := 3
	tableID := RestCreatedTableFlopNext(t, sbPosition, bbPosition)
	restMakeAction(t, tableID.Hex(), domain.Check, bbPosition, secondPlayerToken)

	tableFromDB, err := db.FindTableNoCtx(tableID)
	assert.Nil(t, err)
	assert.EqualValues(t, domain.NewCardStr("3c"), tableFromDB.CommunityCards.Flop.FlopFirst)
}

func TestGetFullTable_NilCardsOfAnotherPlayerForWaitingPlayers(t *testing.T) {
	t.Cleanup(cleanUp)
	table := InsertTable(t)

	restReserveSeat(t, table.ID.Hex(), 0)
	restBuyIn(t, table.ID.Hex(), 0)

	gotTable := fetchTable(t, table.ID)
	assert.Nil(t, gotTable.Seats[0].Player.Cards)
}

func TestGetFullTable_EmptyArrayForComCards_OnEmptyTable(t *testing.T) {
	t.Cleanup(cleanUp)
	table := InsertTable(t)

	gotTable := fetchTable(t, table.ID)

	assert.Len(t, *gotTable.CommunityCards, 0)
}

func TestGetFullTable_ReturnPlayerTimeoutAt(t *testing.T) {
	t.Cleanup(cleanUp)
	tableID := RestCreatedTableWithStartedGame(t, 0, 1)

	gotTable := fetchTable(t, tableID)

	timeoutAt := gotTable.Seats[0].Player.TimeoutAt
	assert.NotNil(t, timeoutAt)
	assert.NotZero(t, timeoutAt)
}

func TestGetFullTable_ForNonAuthUser(t *testing.T) {
	t.Cleanup(cleanUp)
	tableID := RestCreatedTableWithStartedGame(t, 0, 1)

	fetchTable(t, tableID, "")
}

func fetchTable(t *testing.T, id primitive.ObjectID, token ...string) model.Table {
	rr := testRouter.Request(t, http.MethodGet, "/tables/"+id.Hex(), nil, extractAuthHeaders(token))
	assert.EqualValues(t, http.StatusOK, rr.Code)
	var gotTable model.Table
	err := json.Unmarshal(rr.Body.Bytes(), &gotTable)
	assert.Nil(t, err)
	return gotTable
}
