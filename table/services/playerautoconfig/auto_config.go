package playerautoconfig

import (
	"context"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/model"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func SetAll(ctx context.Context, iden authid.Identity, playerAutoConfig *model.PlayerAutoConfig) error {
	return applyAutoConfig(ctx, iden, func(pac *domain.PlayerAutoConfig) {
		pac.Muck = playerAutoConfig.Muck
		pac.TopUp = playerAutoConfig.TopUp
		pac.ReBuy = playerAutoConfig.ReBuy
	})
}

func SetAutoMuck(ctx context.Context, iden authid.Identity, muck bool) error {
	return applyAutoConfig(ctx, iden, func(pac *domain.PlayerAutoConfig) {
		pac.Muck = muck
	})
}

func SetAutoTopUp(ctx context.Context, iden authid.Identity, topUp bool) error {
	return applyAutoConfig(ctx, iden, func(pac *domain.PlayerAutoConfig) {
		pac.TopUp = topUp
	})
}

func SetAutoReBuy(ctx context.Context, iden authid.Identity, reBuy bool) error {
	return applyAutoConfig(ctx, iden, func(pac *domain.PlayerAutoConfig) {
		pac.ReBuy = reBuy
	})
}

func applyAutoConfig(ctx context.Context, iden authid.Identity, apply func(*domain.PlayerAutoConfig)) error {
	pac, err := FindAutoConfig(ctx, iden)
	if err != nil {
		return err
	}
	apply(pac)
	return db.UpsertPlayerAutoConfig(ctx, pac)
}

func FindAutoConfig(ctx context.Context, iden authid.Identity) (*domain.PlayerAutoConfig, error) {
	pac, err := db.FindPlayerAutoConfig(ctx, iden.UserId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			pac := domain.NewPlayerAutoConfig(iden.UserId)
			return &pac, nil
		}
		return nil, err
	}
	return pac, nil
}

func FindAutoConfigNoCtx(iden authid.Identity) (*domain.PlayerAutoConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return FindAutoConfig(ctx, iden)
}
