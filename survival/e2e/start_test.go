package e2e

import (
	"context"
	"github.com/glossd/pokergloss/survival/db"
	"github.com/glossd/pokergloss/survival/domain"
	"github.com/glossd/pokergloss/survival/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestStartSurvival(t *testing.T) {
	t.Cleanup(cleanUpDB)
	ctx := context.Background()
	assert.Nil(t, db.CardIncTwoTickets(ctx, defaultIdentity.UserId))
	post := testRouter.POST(t, "/start", "", authHeaders(defaultToken))
	assert.EqualValues(t, http.StatusCreated, post.Code, post.Body.String())

	_, err := db.Find(ctx, defaultIdentity.UserId)
	assert.Nil(t, err)

	post2 := testRouter.POST(t, "/start", "", authHeaders(defaultToken))
	assert.EqualValues(t, http.StatusOK, post2.Code, post2.Body.String())
}

func TestStartSurvivalIdle(t *testing.T) {
	t.Cleanup(cleanUpDB)
	ctx := context.Background()
	post := testRouter.POST(t, "/start-idle", "", authHeaders(defaultToken))
	assert.EqualValues(t, http.StatusCreated, post.Code, post.Body.String())

	_, err := db.Find(ctx, defaultIdentity.UserId)
	assert.Nil(t, err)
}

func TestSurvivalWithoutSpinning(t *testing.T) {
	t.Cleanup(cleanUpDB)
	ctx := context.Background()
	assert.Nil(t, db.CardIncTwoTickets(ctx, defaultIdentity.UserId))

	res, err := service.Start(ctx, defaultIdentity, domain.Params{})
	assert.Nil(t, err)

	assert.Nil(t, service.EndLevel(ctx, defaultIdentity.UserId, res.TableID, false))
	assert.Nil(t, service.EndLevel(ctx, defaultIdentity.UserId, res.TableID, true))

	res, err = service.Start(ctx, defaultIdentity, domain.Params{})
	assert.Nil(t, err)
	assert.False(t, res.AlreadyStarted)
}
