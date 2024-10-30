package e2e

import (
	"context"
	"github.com/glossd/pokergloss/profile/db"
	"github.com/glossd/pokergloss/profile/web/client/authclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindUsersByUserIDs(t *testing.T) {
	t.Cleanup(cleanUpDB)
	authclient.AuthClient = &authclient.MockAuthClient{UserID: defaultIdentity.UserId}
	createUserFull(t, "test1", "test1@gmail.com", "123456")
	authclient.AuthClient = &authclient.MockAuthClient{UserID: secondIdentity.UserId}
	createUserFull(t, "test2", "test2@gmail.com", "123456")

	profiles, err := db.FindAllProfilesByUserIDs(context.Background(), []string{defaultIdentity.UserId, secondIdentity.UserId})
	assert.Nil(t, err)
	assert.EqualValues(t, 2, len(profiles))
	assert.EqualValues(t, "test1", profiles[0].Username)
	assert.EqualValues(t, "test2", profiles[1].Username)
}
