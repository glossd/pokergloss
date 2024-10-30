package e2e

import (
	"github.com/glossd/pokergloss/bank/db"
	"github.com/glossd/pokergloss/bank/services/fee"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInactionFee(t *testing.T) {
	t.Cleanup(cleanUpDB)
	makeDeposit(t, 200000)
	b, err := db.FindBalanceNoCtx(defaultIdentity.UserId)
	assert.Nil(t, err)
	b.UpdatedAt = time.Now().AddDate(0, 0, -31).UnixNano() / 1e6
	assert.Nil(t, db.UpdateBalanceNoCtx(b))

	fee.RunInactionFee()

	b, err = db.FindBalanceNoCtx(defaultIdentity.UserId)
	assert.Nil(t, err)
	assert.EqualValues(t, 190000, b.Chips)
}

func TestInactionFee_NoFeeForUnderdue(t *testing.T) {
	t.Cleanup(cleanUpDB)
	makeDeposit(t, 200000)
	b, err := db.FindBalanceNoCtx(defaultIdentity.UserId)
	assert.Nil(t, err)
	b.UpdatedAt = time.Now().AddDate(0, 0, -29).UnixNano() / 1e6
	assert.Nil(t, db.UpdateBalanceNoCtx(b))

	fee.RunInactionFee()

	b, err = db.FindBalanceNoCtx(defaultIdentity.UserId)
	assert.Nil(t, err)
	assert.EqualValues(t, 200000, b.Chips)
}

func TestInactionFee_NoFeeForLessThan50000(t *testing.T) {
	t.Cleanup(cleanUpDB)
	makeDeposit(t, 49000)
	b, err := db.FindBalanceNoCtx(defaultIdentity.UserId)
	assert.Nil(t, err)
	b.UpdatedAt = time.Now().AddDate(0, 0, -31).UnixNano() / 1e6
	assert.Nil(t, db.UpdateBalanceNoCtx(b))

	fee.RunInactionFee()

	b, err = db.FindBalanceNoCtx(defaultIdentity.UserId)
	assert.Nil(t, err)
	assert.EqualValues(t, 49000, b.Chips)
}
