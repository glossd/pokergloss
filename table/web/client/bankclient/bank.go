package bankclient

import (
	"context"
	"fmt"
	"github.com/glossd/pokergloss/bank/web/grpc"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gogrpc/grpcbank"
	"github.com/glossd/pokergloss/gomq/mqbank"
	log "github.com/sirupsen/logrus"
	"time"
)

var ErrBankUnavailable = fmt.Errorf("bank is unavailable, try make a deposit later")
var ErrNotEnoughChips = fmt.Errorf("not enough chips")

const DefaultTimeout = time.Second

func WithdrawNoCtx(chips int64, userId, description string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return Withdraw(ctx, chips, userId, description)
}

func Withdraw(ctx context.Context, chips int64, userId, description string) error {
	if chips == 0 {
		return nil
	}
	if conf.IsProd() {
		request := grpcbank.WithdrawRequest{Chips: chips, UserId: userId, Description: description}
		res, err := grpc.Withdraw(ctx, &request)
		if err != nil {
			return err
		}
		if res.Status == grpcbank.WithdrawResponse_NOT_ENOUGH_CHIPS {
			return ErrNotEnoughChips
		}
	}

	return nil
}

func GetBalanceWithRank(userId string) (balance int64, rank int64, err error) {
	if conf.IsProd() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := grpc.GetBalanceWithRank(ctx, &grpcbank.GetBalanceWithRankRequest{UserId: userId})
		if err != nil {
			return 0, 0, err
		}
		return r.Chips, r.Rank, nil
	}
	return 0, 0, nil
}

func Deposit(chips int64, userId, description string) error {
	if chips == 0 {
		return nil
	}
	if conf.IsProd() {
		request := &mqbank.MultiDeposit{
			Service: "table",
			Deposits: toUserDeposits([]UserDeposit{
				{UserID: userId, Amount: chips, Description: description}},
			)}
		err := mqbank.PublishMultiDepositAsync(request)
		if err != nil {
			log.Errorf("Couldn't send a deposit to pubsub, userId=%s, chips:%d", userId, chips)
		}
		return err
	}
	return nil
}

type UserDeposit struct {
	UserID      string
	Amount      int64
	Description string
}

func DepositSync(ctx context.Context, deposit UserDeposit) error {
	if !conf.IsProd() {
		return nil
	}
	err := mqbank.PublishMultiDeposit(ctx, &mqbank.MultiDeposit{
		Service:  "table",
		Deposits: []*mqbank.UserDeposit{toUserDeposit(deposit)},
	})
	if err != nil {
		log.Errorf("Failed to send deposit: %s", err)
	}
	return err
}

func MultiDeposit(ctx context.Context, deposists []UserDeposit) error {
	if !conf.IsProd() {
		return nil
	}
	err := mqbank.PublishMultiDeposit(ctx, &mqbank.MultiDeposit{
		Deposits: toUserDeposits(deposists),
		Service:  "table",
	})
	if err != nil {
		log.Errorf("Failed to send multi deposit: %s", err)
	}
	return err
}

func MultiDepositAsync(deposists []UserDeposit) error {
	if !conf.IsProd() {
		return nil
	}
	return mqbank.PublishMultiDepositAsync(&mqbank.MultiDeposit{
		Deposits: toUserDeposits(deposists),
		Service:  "table",
	})
}

func toUserDeposits(d []UserDeposit) []*mqbank.UserDeposit {
	result := make([]*mqbank.UserDeposit, 0, len(d))
	for _, deposit := range d {
		if deposit.Amount == 0 {
			continue
		}
		result = append(result, toUserDeposit(deposit))
	}
	return result
}

func toUserDeposit(d UserDeposit) *mqbank.UserDeposit {
	return &mqbank.UserDeposit{
		Amount:      d.Amount,
		UserId:      d.UserID,
		Description: d.Description,
	}
}
