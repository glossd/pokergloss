package e2e

import (
	"context"
	"github.com/glossd/pokergloss/profile/db"
	"github.com/glossd/pokergloss/profile/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckUsernameDoesNotExist(t *testing.T) {
	t.Cleanup(cleanUpDB)
	rr := testRouter.Request(t, "GET", "/users/check-username?username=goose", nil, nil)
	assert.EqualValues(t, 200, rr.Code)
	assert.EqualValues(t, `{"message":"username is unique"}`, rr.Body.String())
}

func TestCheckUsernameExists(t *testing.T) {
	t.Cleanup(cleanUpDB)
	err := db.UpsertProfile(context.TODO(), domain.NewProfile("goose", "userID", defaultTechInfo))
	assert.Nil(t, err)
	rr := testRouter.Request(t, "GET", "/users/check-username?username=goose", nil, nil)
	assert.EqualValues(t, 200, rr.Code)
	assert.EqualValues(t, `{"message":"username already exists"}`, rr.Body.String())
}
