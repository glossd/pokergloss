package e2e

import (
	"github.com/glossd/pokergloss/achievement/service"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"testing"
)

func TestStatistics(t *testing.T) {
	t.Cleanup(cleanUp)
	service.UpdateStatistics(oneWinnerTwoPlayers())
	rr := testRouter.Request(t, http.MethodGet, "/statistics/me", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	assert.EqualValues(t, 1, gjson.Get(body, "gameCount").Int())
	assert.EqualValues(t, 100, gjson.Get(body, "winPercent").Int())
	assert.EqualValues(t, 100, gjson.Get(body, "allInPercent").Int())
	assert.EqualValues(t, 0, gjson.Get(body, "foldPercent").Int())

	service.UpdateStatistics(oneDefaultLosesTwoPlayers())
	rr = testRouter.Request(t, http.MethodGet, "/statistics/me", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	body = rr.Body.String()
	assert.EqualValues(t, 2, gjson.Get(body, "gameCount").Int())
	assert.EqualValues(t, 50, gjson.Get(body, "winPercent").Int())
	assert.EqualValues(t, 50, gjson.Get(body, "allInPercent").Int())
	assert.EqualValues(t, 50, gjson.Get(body, "foldPercent").Int())
}
