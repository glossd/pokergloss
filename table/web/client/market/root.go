package market

import (
	"context"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gogrpc/grpcmarket"
	"github.com/glossd/pokergloss/market/web/grpc"
	"time"
)

type SelectedItem struct {
	ID            string
	CoinsDayPrice int64
}

func GetSelectedItemID(userID string) (SelectedItem, error) {
	if !conf.IsProd() {
		return SelectedItem{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r := &grpcmarket.GetUserSelectedItemRequest{UserId: userID}
	res, err := grpc.GetUserSelectedItem(ctx, r)
	if err != nil {
		return SelectedItem{}, err
	}

	return SelectedItem{ID: res.ItemId, CoinsDayPrice: res.CoinsDayPrice}, nil
}
