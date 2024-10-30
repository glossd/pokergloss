package e2e

import (
	"context"
	"github.com/glossd/pokergloss/bank/db"
	"github.com/glossd/pokergloss/bank/domain"
	"github.com/glossd/pokergloss/gogrpc/grpcbank"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestWithdraw(t *testing.T) {
	t.Cleanup(cleanUpDB)
	t.Run("chips", func(t *testing.T) {
		makeDeposit(t, 300)
		res := makeWithdraw(t, 200)
		assert.EqualValues(t, res.GetStatus(), grpcbank.WithdrawResponse_OK)
		assert.EqualValues(t, 100, findBalance(t).Chips)
	})
	t.Run("coins", func(t *testing.T) {
		makeDepositCoins(t, 300)
		res := makeWithdrawCoins(t, 200)
		assert.EqualValues(t, res.GetStatus(), grpcbank.WithdrawCoinsResponse_OK)
		assert.EqualValues(t, 100, findBalance(t).Coins)
	})
}

func TestWithdraw_NotEnough(t *testing.T) {
	t.Cleanup(cleanUpDB)
	t.Run("chips", func(t *testing.T) {
		makeDeposit(t, 100)
		res := makeWithdraw(t, 200)
		assert.EqualValues(t, res.GetStatus(), grpcbank.WithdrawResponse_NOT_ENOUGH_CHIPS)
	})
	t.Run("coins", func(t *testing.T) {
		makeDepositCoins(t, 100)
		res := makeWithdrawCoins(t, 200)
		assert.EqualValues(t, res.GetStatus(), grpcbank.WithdrawCoinsResponse_NOT_ENOUGH_COINS)
	})
}

func TestMakeWithdraw_WithoutDeposit(t *testing.T) {
	t.Cleanup(cleanUpDB)

	t.Run("chips", func(t *testing.T) {
		res := makeWithdraw(t, 200)
		assert.EqualValues(t, res.GetStatus(), grpcbank.WithdrawResponse_NOT_ENOUGH_CHIPS)
	})

	t.Run("coins", func(t *testing.T) {
		res := makeWithdrawCoins(t, 200)
		assert.EqualValues(t, res.GetStatus(), grpcbank.WithdrawCoinsResponse_NOT_ENOUGH_COINS)
	})
}

func makeWithdraw(t *testing.T, chips int64) *grpcbank.WithdrawResponse {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	assert.Nil(t, err)
	defer conn.Close()
	client := grpcbank.NewBankServiceClient(conn)
	res, err := client.Withdraw(ctx, &grpcbank.WithdrawRequest{
		Chips:       chips,
		Description: "welcome bonus",
		UserId:      defaultIdentity.UserId,
	})
	assert.Nil(t, err)
	return res
}

func makeWithdrawCoins(t *testing.T, coins int64) *grpcbank.WithdrawCoinsResponse {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	assert.Nil(t, err)
	defer conn.Close()
	client := grpcbank.NewBankServiceClient(conn)
	res, err := client.WithdrawCoins(ctx, &grpcbank.WithdrawCoinsRequest{
		Coins:       coins,
		Description: "welcome bonus",
		UserId:      defaultIdentity.UserId,
	})
	assert.Nil(t, err)
	return res
}

func findBalance(t *testing.T) *domain.Balance {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	b, err := db.FindBalance(ctx, defaultIdentity.UserId)
	assert.Nil(t, err)
	return b
}
