package bank

import (
	"context"
	"fmt"
	"github.com/glossd/pokergloss/bank/web/grpc"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gogrpc/grpcbank"
	"github.com/glossd/pokergloss/gomq/mqbank"
	log "github.com/sirupsen/logrus"
)

var ErrNotEnoughChips = fmt.Errorf("not enough chips")
var ErrNotEnoughCoins = fmt.Errorf("not enough coins")

func Withdraw(ctx context.Context, chips int64, userId, description string) error {
	if conf.IsProd() {
		request := grpcbank.WithdrawRequest{Chips: chips, UserId: userId, Description: description}
		res, err := grpc.Withdraw(ctx, &request)
		if err != nil {
			log.Errorf("bank.Withdraw: failed rpc request Withdraw: %s", err)
			return err
		}
		if res.Status == grpcbank.WithdrawResponse_NOT_ENOUGH_CHIPS {
			return ErrNotEnoughChips
		}
	}

	return nil
}

func WithdrawCoins(ctx context.Context, chips int64, userId, description string) error {
	if conf.IsProd() {
		request := grpcbank.WithdrawCoinsRequest{Coins: chips, UserId: userId, Reason: "market", Description: description}
		res, err := grpc.WithdrawCoins(ctx, &request)
		if err != nil {
			log.Errorf("bank.WithdrawCoins: failed rpc request WithdrawCoins: %s", err)
			return err
		}
		if res.Status == grpcbank.WithdrawCoinsResponse_NOT_ENOUGH_COINS {
			return ErrNotEnoughCoins
		}
	}

	return nil
}

func Deposit(chips int64, userId, description string) {
	deposit(chips, userId, description, mqbank.CurrencyType_CHIPS)
}

func DepositCoins(coins int64, userId, description string) {
	deposit(coins, userId, description, mqbank.CurrencyType_COINS)
}

func deposit(value int64, userId string, description string, currency mqbank.CurrencyType) {
	if value == 0 {
		return
	}
	if conf.IsProd() {
		request := mqbank.DepositRequest{Chips: value, UserId: userId, Description: description, Reason: "market", CurrencyType: currency}
		err := mqbank.Deposit(&request)
		if err != nil {
			if currency == mqbank.CurrencyType_COINS {
				log.Errorf("Couldn't send a coins deposit to pubsub, userId=%s, value=%d", userId, value)
			} else {
				log.Errorf("Couldn't send a chips deposit to pubsub, userId=%s, value=%d", userId, value)
			}
		}
	}
}
