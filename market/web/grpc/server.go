package grpc

import (
	"context"
	"github.com/glossd/pokergloss/gogrpc/grpcmarket"
	"github.com/glossd/pokergloss/market/domain"
	"github.com/glossd/pokergloss/market/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetUserSelectedItem(ctx context.Context, r *grpcmarket.GetUserSelectedItemRequest) (*grpcmarket.GetUserSelectedItemResponse, error) {
	item, err := service.GetSelectedItem(ctx, r.GetUserId())
	if err != nil {
		if _, ok := err.(*domain.Err); ok {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		if _, ok := err.(*service.Err); ok {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	return &grpcmarket.GetUserSelectedItemResponse{ItemId: string(item.ItemID), CoinsDayPrice: domain.ItemCoinsDayPrice(item.ItemID)}, nil
}
