package e2e

import (
	"github.com/glossd/pokergloss/profile/web/client/authclient"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	t.Cleanup(cleanUpDB)
	authclient.AuthClient = &authclient.MockAuthClient{}
	createUser(t)

	rr := testRouter.Request(t, http.MethodGet, "/profiles/den", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, "den", gjson.Get(rr.Body.String(), "username").String())
}

func TestSearchProfiles(t *testing.T) {
	t.Cleanup(cleanUpDB)
	authclient.AuthClient = &authclient.MockAuthClient{}
	createUser(t)

	rr := testRouter.Request(t, http.MethodGet, "/profiles/de/search", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 1, gjson.Get(rr.Body.String(), "#").Int())
	assert.EqualValues(t, "den", gjson.Get(rr.Body.String(), "0.username").String())
}
