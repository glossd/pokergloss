package achievement

import (
	"context"
	"github.com/glossd/pokergloss/achievement/web/grpc"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gogrpc/grpcachieve"
	"time"
)

func GetLevel(userID string) (int64, error) {
	if !conf.IsProd() {
		return 1, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	r := &grpcachieve.GetExpRequest{UserId: userID}
	exp, err := grpc.GetExp(ctx, r)
	if err != nil {
		return 0, err
	}

	return exp.Level, nil
}
