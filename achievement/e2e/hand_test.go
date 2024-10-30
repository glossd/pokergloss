package e2e

import (
	"github.com/glossd/pokergloss/achievement/db"
	"github.com/glossd/pokergloss/achievement/domain"
	"github.com/glossd/pokergloss/achievement/service"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"testing"
)

func TestSaveAchievement(t *testing.T) {
	t.Cleanup(cleanUp)
	gameEnd := oneWinnerTwoPlayers()
	service.UpdateAchievementStoreNoCtx(gameEnd)
	as, err := db.FindAchievementStoreNoCtx(defaultIdentity.UserId)
	assert.Nil(t, err)
	hc := as.HandsCounter
	assert.EqualValues(t, 7, len(hc.Hands))
	assert.EqualValues(t, 1, hc.Hands[domain.S].Count)
	assert.EqualValues(t, 1, hc.Hands[domain.S].Level)

	rr := testRouter.Request(t, http.MethodGet, "/achievements/me", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	assert.EqualValues(t, "Straight", gjson.Get(body, "0.name").String())
	assert.EqualValues(t, 1, gjson.Get(body, "0.count").Int())
	assert.EqualValues(t, 5, gjson.Get(body, "0.nextLevelCount").Int())
}
