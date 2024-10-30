package grpc

import (
	"context"
	"errors"
	"github.com/glossd/pokergloss/bank/domain"
	"github.com/glossd/pokergloss/bank/services"
	"github.com/glossd/pokergloss/bank/services/balance"
	"github.com/glossd/pokergloss/gogrpc/grpcbank"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Withdraw(ctx context.Context, r *grpcbank.WithdrawRequest) (*grpcbank.WithdrawResponse, error) {
	var reason domain.Reason
	if r.Reason == "" {
		reason = domain.Game
	} else {
		reason = domain.Reason(r.Reason)
	}
	operation, err := domain.NewWithdraw(reason, r.GetChips(), r.GetUserId(), r.GetDescription())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = services.Withdraw(ctx, operation)
	if err != nil {
		if errors.Is(err, domain.ErrNotEnoughChips) {
			return &grpcbank.WithdrawResponse{Status: grpcbank.WithdrawResponse_NOT_ENOUGH_CHIPS}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpcbank.WithdrawResponse{Status: grpcbank.WithdrawResponse_OK}, nil
}

func GetBalance(ctx context.Context, r *grpcbank.GetBalanceRequest) (*grpcbank.GetBalanceResponse, error) {
	b, err := balance.Get(ctx, r.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpcbank.GetBalanceResponse{
		UserId: b.UserID,
		Chips:  b.Chips,
	}, nil
}

func GetBalanceWithRank(ctx context.Context, r *grpcbank.GetBalanceWithRankRequest) (*grpcbank.GetBalanceWithRankResponse, error) {
	b, err := services.GetBalanceWithRank(ctx, r.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpcbank.GetBalanceWithRankResponse{
		UserId: r.UserId,
		Chips:  b.Balance,
		Rank:   b.Rank,
	}, nil
}

func WithdrawCoins(ctx context.Context, r *grpcbank.WithdrawCoinsRequest) (*grpcbank.WithdrawCoinsResponse, error) {
	withdraw, err := domain.NewWithdrawCoins(domain.Reason(r.Reason), r.Coins, r.UserId, r.Description)
	if err != nil {
		return nil, err
	}
	err = services.WithdrawCoins(ctx, withdraw)
	if err != nil {
		if err == domain.ErrNotEnoughCoins {
			return &grpcbank.WithdrawCoinsResponse{Status: grpcbank.WithdrawCoinsResponse_NOT_ENOUGH_COINS}, nil
		}
		return nil, err
	}
	return &grpcbank.WithdrawCoinsResponse{Status: grpcbank.WithdrawCoinsResponse_OK}, nil
}
