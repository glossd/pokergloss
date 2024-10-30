package e2e

import (
	"github.com/glossd/pokergloss/bonus/db"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"testing"
)

func TestDailyBonus(t *testing.T) {
	rr := testRouter.Request(t, http.MethodPut, "/daily-bonus", nil, nil)
	body := rr.Body.String()
	assert.EqualValues(t, http.StatusOK, rr.Code, body)

	assert.EqualValues(t, true, gjson.Get(body, "isBonusPresent").Bool())
	assert.EqualValues(t, 500, gjson.Get(body, "bonus.chips").Int())
	assert.EqualValues(t, 1, gjson.Get(body, "bonus.dayInARow").Int())

	rr = testRouter.Request(t, http.MethodPut, "/daily-bonus", nil, nil)
	body = rr.Body.String()
	assert.EqualValues(t, http.StatusOK, rr.Code, body)
	assert.EqualValues(t, false, gjson.Get(body, "isBonusPresent").Bool())

	db.ResetBonuses()

	rr = testRouter.Request(t, http.MethodPut, "/daily-bonus", nil, nil)
	body = rr.Body.String()
	assert.EqualValues(t, http.StatusOK, rr.Code, body)

	assert.EqualValues(t, true, gjson.Get(body, "isBonusPresent").Bool())
	assert.EqualValues(t, 1207, gjson.Get(body, "bonus.chips").Int())
	assert.EqualValues(t, 2, gjson.Get(body, "bonus.dayInARow").Int())
}
