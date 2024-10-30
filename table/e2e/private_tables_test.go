package e2e

import (
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPrivateTables_ShouldNotAppearList(t *testing.T) {
	t.Cleanup(cleanUp)
	postTable(t)
	postTablePrivate(t)
	body := findCashTable(t)
	assert.EqualValues(t, 1, gjson.Get(body, "#").Int())
}

func postTablePrivate(t *testing.T) *httptest.ResponseRecorder {
	postBody := `{"name":"my private table", "size":9, "bigBlind":2, "isPrivate": true}`
	rr := testRouter.Request(t, http.MethodPost, "/tables", &postBody, nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	return rr
}

func findCashTable(t *testing.T) (body string) {
	rr := testRouter.Request(t, http.MethodGet, "/tables", nil, nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	return rr.Body.String()
}

