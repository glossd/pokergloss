package e2e

import (
	"context"
	"github.com/glossd/pokergloss/profile/db"
	"github.com/glossd/pokergloss/profile/domain"
	"github.com/glossd/pokergloss/profile/web/client/authclient"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestSignUp(t *testing.T) {
	t.Cleanup(cleanUpDB)
	authclient.AuthClient = &authclient.MockAuthClient{}
	body := "username=den&email=den@mail.ru&password=123456"
	rr := testRouter.Request(t, http.MethodPost, "/signup", &body, map[string]string{"Content-type": "application/x-www-form-urlencoded"})
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	exists, err := db.ExistsUsername(ctx, "den")
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestFailSignUp_UsernameExists(t *testing.T) {
	t.Cleanup(cleanUpDB)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := db.UpsertProfile(ctx, domain.NewProfile("den", "", defaultTechInfo))
	assert.Nil(t, err)

	authclient.AuthClient = &authclient.MockAuthClient{}
	body := "username=den&email=den@mail.ru&password=123456"
	rr := testRouter.Request(t, http.MethodPost, "/signup", &body, map[string]string{"Content-type": "application/x-www-form-urlencoded"})
	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
	assert.EqualValues(t, `{"message":"username is taken"}`, rr.Body.String())
}
