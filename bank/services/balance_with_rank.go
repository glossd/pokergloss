package services

import (
	"context"
	"github.com/glossd/pokergloss/bank/services/balance"
	"github.com/glossd/pokergloss/bank/services/model"
)

func GetBalanceWithRank(ctx context.Context, userID string) (*model.BalanceWithRank, error) {
	b, err := balance.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if b.NonRatable {
		return &model.BalanceWithRank{
			Balance: b.Chips,
			Coins:   b.Coins,
			Rank:    0,
		}, nil
	}
	r, err := GetRating(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &model.BalanceWithRank{
		Balance: b.Chips,
		Coins:   b.Coins,
		Rank:    r.Rank,
	}, nil
}
