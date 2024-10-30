package e2e

import (
	"context"
	"github.com/glossd/pokergloss/bank/db"
	"github.com/glossd/pokergloss/bank/domain"
	"github.com/glossd/pokergloss/bank/services"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestGetRankWithNoView(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := db.FindRating(ctx, defaultIdentity.UserId)
	assert.Equal(t, err, mongo.ErrNoDocuments)
}

func TestRank(t *testing.T) {
	t.Cleanup(cleanUpDB)

	deposit(t, defaultIdentity.UserId, 1000)

	deposit(t, "richer", 2000)

	db.BuildOppositeRankView()

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	rating, err := db.FindRating(ctx, defaultIdentity.UserId)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, rating.Rank)
	assert.EqualValues(t, 1000, rating.Chips)
	assert.NotNil(t, rating.UpdatedAt)

	deposit(t, defaultIdentity.UserId, 5000)
	db.BuildOppositeRankView()

	rating, err = db.FindRating(ctx, defaultIdentity.UserId)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, rating.Rank)
	assert.EqualValues(t, 6000, rating.Chips)
	assert.NotNil(t, rating.UpdatedAt)
}

func deposit(t *testing.T, userID string, chips int64) {
	_, err := db.FindBalanceNoCtx(userID)
	if err == mongo.ErrNoDocuments {
		assert.Nil(t, services.CreateFirstBalance(context.Background(), domain.ProfileInfo{UserID: userID, Username: userID}))
	}
	deposit, err := domain.NewDeposit(domain.Bonus, chips, userID, "")
	assert.Nil(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = services.Deposit(ctx, deposit)
	assert.Nil(t, err)
}
