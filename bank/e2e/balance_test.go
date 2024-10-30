package e2e

import (
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"testing"
)

func TestGetBalance(t *testing.T) {
	t.Cleanup(cleanUpDB)
	makeDeposit(t, 200)

	rr := testRouter.Request(t, http.MethodGet, "/balance", nil, nil)
	assert.EqualValues(t , http.StatusOK, rr.Code, rr.Body.String())

	assert.EqualValues(t, 200, gjson.Get(rr.Body.String(), "chips").Int())
}

func TestGetBalance_IfNoOperations(t *testing.T) {
	t.Cleanup(cleanUpDB)
	rr := testRouter.Request(t, http.MethodGet, "/balance", nil, nil)
	assert.EqualValues(t , http.StatusOK, rr.Code, rr.Body.String())
	assert.EqualValues(t, 0, gjson.Get(rr.Body.String(), "chips").Int())
}
