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

func TestDeposit(t *testing.T) {
	t.Cleanup(cleanUpDB)
	makeDeposit(t, 200)
	b, err := db.FindBalanceNoCtx(defaultIdentity.UserId)
	assert.Nil(t, err)
	assert.EqualValues(t, 200, b.Chips)
}

func makeDeposit(t *testing.T, chips int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := db.FindBalance(ctx, defaultIdentity.UserId)
	if err == mongo.ErrNoDocuments {
		err := services.CreateFirstBalance(ctx, domain.ProfileInfo{UserID: defaultIdentity.UserId, Username: defaultIdentity.Username})
		assert.Nil(t, err)
	}
	newDeposit, err := domain.NewDeposit(domain.Bonus, chips, defaultIdentity.UserId, "welcome bonus")
	assert.Nil(t, err)
	err = services.Deposit(ctx, newDeposit)
	assert.Nil(t, err)
}

func makeDepositCoins(t *testing.T, coins int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := db.FindBalance(ctx, defaultIdentity.UserId)
	if err == mongo.ErrNoDocuments {
		err := services.CreateFirstBalance(ctx, domain.ProfileInfo{UserID: defaultIdentity.UserId, Username: defaultIdentity.Username})
		assert.Nil(t, err)
	}
	newDeposit, err := domain.NewDepositCoins(domain.Bonus, coins, defaultIdentity.UserId, "welcome bonus")
	assert.Nil(t, err)
	err = services.DepositCoins(ctx, newDeposit)
	assert.Nil(t, err)
}
