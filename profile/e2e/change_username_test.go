package e2e

import (
	"github.com/glossd/pokergloss/profile/db"
	"github.com/glossd/pokergloss/profile/web/client/authclient"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestChangeUsernameForNonExistingProfile(t *testing.T) {
	inputBody := `{"username":"Se"}`
	rr := authTestRouter.Request(t, http.MethodPut, "/users/me/username/change", &inputBody, nil)
	assert.EqualValues(t, http.StatusBadRequest, rr.Code, rr.Body.String())
}

func TestChangeUsernameOnExisting(t *testing.T) {
	t.Cleanup(cleanUpDB)
	authclient.AuthClient = &authclient.MockAuthClient{}
	createUser(t)
	inputBody := `{"username":"den"}`
	rr := authTestRouter.Request(t, http.MethodPut, "/users/me/username/change", &inputBody, authHeaders(secondToken))
	assert.EqualValues(t, http.StatusBadRequest, rr.Code, rr.Body.String())
}

func TestChangeUsername(t *testing.T) {
	t.Cleanup(cleanUpDB)
	authclient.AuthClient = &authclient.MockAuthClient{UserID: defaultIdentity.UserId}
	createUserWithUsername(t, defaultIdentity.Username)
	inputBody := `{"username":"den2"}`
	rr := authTestRouter.Request(t, http.MethodPut, "/users/me/username/change", &inputBody, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())

	p, err := db.FindProfileNoCtx("den2")
	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(p.OldUsernames))
	assert.EqualValues(t, defaultIdentity.Username, p.OldUsernames[0].Username)
}
