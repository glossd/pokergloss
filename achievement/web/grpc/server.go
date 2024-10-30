package grpc

import (
	"context"
	"github.com/glossd/pokergloss/achievement/service"
	"github.com/glossd/pokergloss/gogrpc/grpcachieve"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetExp(ctx context.Context, r *grpcachieve.GetExpRequest) (*grpcachieve.GetExpResponse, error) {
	exp, err := service.GetUserExp(ctx, r.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &grpcachieve.GetExpResponse{
		UserId: exp.UserID,
		Points: exp.Points,
		Level:  int64(exp.Level),
	}, nil
}
