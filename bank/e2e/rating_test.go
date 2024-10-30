package e2e

import (
	"github.com/glossd/pokergloss/bank/db"
	"github.com/glossd/pokergloss/bank/services/ranker"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"testing"
	"time"
)

func TestGetRatings(t *testing.T) {
	t.Cleanup(cleanUpDB)
	deposit(t, "first", 5000)
	deposit(t, "second", 1000)
	deposit(t, "third", 500)
	db.BuildOppositeRankView()

	rr := testRouter.Request(t, http.MethodGet, "/ratings/page?pageSize=20&pageNumber=1", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)

	body := rr.Body.String()
	assert.EqualValues(t, 3, gjson.Get(body, "ratings.#").Int())
	assert.EqualValues(t, 5000, gjson.Get(body, "ratings.0.chips").Int())
	assert.EqualValues(t, "first", gjson.Get(body, "ratings.0.username").String())
}

func TestGetRatings_NoPageNumber(t *testing.T) {
	t.Cleanup(cleanUpDB)
	deposit(t, "first", 5000)
	deposit(t, "second", 1000)
	deposit(t, "third", 500)
	db.BuildOppositeRankView()

	rr := testRouter.Request(t, http.MethodGet, "/ratings/page?pageSize=20", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)

	body := rr.Body.String()
	assert.EqualValues(t, 3, gjson.Get(body, "ratings.#").Int())
	assert.EqualValues(t, 5000, gjson.Get(body, "ratings.0.chips").Int())
}

func TestGetRatings_NoPageNumber_ExistingUser(t *testing.T) {
	t.Cleanup(cleanUpDB)
	deposit(t, "first", 5000)
	deposit(t, "second", 1000)
	deposit(t, "third", 500)
	deposit(t, "2fR5MYyjqMSMgzLtdkdrKZkc18t1", 100)
	db.BuildOppositeRankView()

	rr := testRouter.Request(t, http.MethodGet, "/ratings/page?pageSize=2", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)

	body := rr.Body.String()
	assert.EqualValues(t, 2, gjson.Get(body, "ratings.#").Int())
	assert.EqualValues(t, 100, gjson.Get(body, "ratings.1.chips").Int())
}

func TestGetTopRatings(t *testing.T) {
	t.Cleanup(cleanUpDB)
	deposit(t, "first", 5000)
	deposit(t, "second", 1000)
	deposit(t, "third", 500)
	deposit(t, "fourth", 250)
	deposit(t, "fourth", 250)

	tChan := make(chan time.Time)
	ticker := &time.Ticker{C: tChan}
	go func() {
		tChan <- time.Now()
		close(tChan)
	}()
	ranker.RunRanker(ticker)

	rr := testRouter.Request(t, http.MethodGet, "/ratings/page?pageSize=10&pageNumber=1", nil, map[string]string{"Authorization": "blah"})
	assert.EqualValues(t, http.StatusOK, rr.Code)

	body := rr.Body.String()
	assert.EqualValues(t, 4, gjson.Get(body, "ratings.#").Int())

	assert.EqualValues(t, "fourth", gjson.Get(body, "ratings.2.userId").String())
	assert.EqualValues(t, 3, gjson.Get(body, "ratings.2.rank").Int())
	assert.EqualValues(t, 500, gjson.Get(body, "ratings.2.chips").Int())

	assert.EqualValues(t, "third", gjson.Get(body, "ratings.3.userId").String())
}
