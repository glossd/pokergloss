package e2e

import (
	"github.com/glossd/pokergloss/achievement/service"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"testing"
)

func TestGetAchievements(t *testing.T) {
	t.Cleanup(cleanUp)
	rr := testRouter.Request(t, http.MethodGet, "/achievements/me", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	assert.EqualValues(t, 12, gjson.Get(body, "#").Int())
	assert.EqualValues(t, "Two Pair", gjson.Get(body, "0.name").String())
}

func TestGetAchievementsSorted(t *testing.T) {
	t.Cleanup(cleanUp)
	gameEnd := oneWinnerTwoPlayers()
	service.UpdateAchievementStoreNoCtx(gameEnd)
	rr := testRouter.Request(t, http.MethodGet, "/achievements/me", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	assert.EqualValues(t, 12, gjson.Get(body, "#").Int())
	assert.EqualValues(t, "Straight", gjson.Get(body, "0.name").String())
	assert.EqualValues(t, 1, gjson.Get(body, "0.level").Int())
	assert.EqualValues(t, "Two Pair", gjson.Get(body, "1.name").String())
	assert.EqualValues(t, 0, gjson.Get(body, "1.level").Int())
}
