package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func restReserveSeat(t *testing.T, tableID string, position int, token ...string) {
	rr := testRouter.Request(t, "POST", "/tables/"+tableID+"/seats/"+strconv.Itoa(position)+"/reserve", nil, extractAuthHeaders(token))
	assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String())
}

func restCancelSeatReservation(t *testing.T, tableID string, position int, token ...string) {
	rr := testRouter.Request(t, "DELETE", fmt.Sprintf("/tables/%s/seats/%d/cancel", tableID, position), nil, extractAuthHeaders(token))
	assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String())
}

func restBuyIn(t *testing.T, tableID string, position int, token ...string) {
	postBody := fmt.Sprintf(`{"stack":250}`)
	rr := testRouter.Request(t, "PUT", "/tables/"+tableID+"/seats/"+strconv.Itoa(position)+"/buy-in", &postBody, extractAuthHeaders(token))
	assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String())
}

func restBuyInStatus(t *testing.T, tableID string, position int, code int, token ...string) {
	postBody := fmt.Sprintf(`{"stack":250}`)
	rr := testRouter.Request(t, "PUT", "/tables/"+tableID+"/seats/"+strconv.Itoa(position)+"/buy-in", &postBody, extractAuthHeaders(token))
	assert.Equal(t, code, rr.Code, rr.Body.String())
}

func restAddChips(t *testing.T, tableID string, position int, token ...string) {
	postBody := fmt.Sprintf(`{"chips":10}`)
	rr := testRouter.Request(t, "PUT", "/tables/"+tableID+"/seats/"+strconv.Itoa(position)+"/add-chips", &postBody, extractAuthHeaders(token))
	assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String())
}

func restSitBack(t *testing.T, tableID string, position int, token ...string) {
	rr := testRouter.Request(t, "PUT", fmt.Sprintf("/tables/%s/seats/%d/sit-back", tableID, position), nil, extractAuthHeaders(token))
	assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String())
}

func restMakeAction(t *testing.T, tableID string, action domain.ActionType, position int, token ...string) {
	body := fmt.Sprintf(`{"chips":0}`)
	rr := testRouter.Request(t, "PUT", "/tables/"+tableID+"/seats/"+strconv.Itoa(position)+"/actions/"+string(action), &body, extractAuthHeaders(token))
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())
}

// raise or bet
func restMakeBetAction(t *testing.T, tableID string, action domain.ActionType, chips int, position int, token ...string) {
	body := fmt.Sprintf(`{"chips":%d}`, chips)
	rr := testRouter.Request(t, "PUT", "/tables/"+tableID+"/seats/"+strconv.Itoa(position)+"/actions/"+string(action), &body, extractAuthHeaders(token))
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())
}

func restMakeActionStatus(t *testing.T, tableID string, action domain.ActionType, position int, code int, token ...string) *httptest.ResponseRecorder {
	body := fmt.Sprintf(`{"chips":0}`)
	rr := testRouter.Request(t, "PUT", "/tables/"+tableID+"/seats/"+strconv.Itoa(position)+"/actions/"+string(action), &body, extractAuthHeaders(token))
	assert.EqualValues(t, code, rr.Code, rr.Body.String())
	return rr
}

func restPlayerStand(t *testing.T, tableID primitive.ObjectID, position int, token ...string) {
	body := fmt.Sprintf(`{"position":%d}`, position) // todo remove
	rr := testRouter.Request(t, http.MethodDelete, urlTable(tableID.Hex(), position, "stand"), &body, extractAuthHeaders(token))
	assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String())
}

func restPutIntent(t *testing.T, tableID string, position int, intent *model.Intent, token ...string) {
	body := fmt.Sprintf(`{"intent": {"type":"%s", "chips": %d}}`, intent.Type, intent.Chips)
	rr := testRouter.Request(t, "PUT", fmt.Sprintf("/tables/%s/seats/%d/intent", tableID, position), &body, extractAuthHeaders(token))
	assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String())
}

func restDeleteIntent(t *testing.T, tableID string, position int, token ...string) {
	rr := testRouter.Request(t, "DELETE", fmt.Sprintf("/tables/%s/seats/%d/intent", tableID, position), nil, extractAuthHeaders(token))
	assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String())
}

func restMakeShowDownAction(t *testing.T, tableID string, action domain.ShowDownActionType, position int, token ...string) {
	rr := testRouter.Request(t, "PUT", urlTable(tableID, position, "showdown-actions/"+string(action)), nil, extractAuthHeaders(token))
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())
}

func restSetAutoMuckConfig(t *testing.T, tableID string, position int, autoMuck bool, token ...string) {
	body := fmt.Sprintf(`{"autoMuck":%t}`, autoMuck)
	rr := testRouter.Request(t, "PUT", urlTable(tableID, position, "configs/auto-muck"), &body, extractAuthHeaders(token))
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())
}

func restSitngoRegister(t *testing.T, id string, position int, token ...string) {
	body := fmt.Sprintf(`{"position":%d}`, position)
	rr := testRouter.Request(t, "PUT", urlSitngo(id, "register"), &body, extractAuthHeaders(token))
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())
}

func restGetFullTable(t *testing.T, id string, token ...string) *model.Table {
	rr := testRouter.Request(t, http.MethodGet, "/tables/"+id, nil, extractAuthHeaders(token))
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())
	var table model.Table
	err := json.Unmarshal(rr.Body.Bytes(), &table)
	assert.Nil(t, err)
	return &table
}

func urlTable(tableId string, position int, rest string) string {
	return fmt.Sprintf("/tables/%s/seats/%d/%s", tableId, position, rest)
}

func urlSitngo(tableId string, rest string) string {
	return fmt.Sprintf("/sit-n-go/lobbies/%s/%s", tableId, rest)
}

func extractAuthHeaders(token []string) map[string]string {
	if len(token) > 0 {
		return authHeaders(token[0])
	}
	return nil
}
