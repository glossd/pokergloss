package db

import (
	"context"
	"github.com/glossd/pokergloss/table/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var True = true

func UpsertPlayerAutoConfig(ctx context.Context, pac *domain.PlayerAutoConfig) error {
	_, err := ColPlayerAutoConfig().ReplaceOne(ctx, filterID(pac.UserID), pac, &options.ReplaceOptions{Upsert: &True})
	if err != nil {
		return err
	}
	return nil
}

func FindPlayerAutoConfig(ctx context.Context, userID string) (*domain.PlayerAutoConfig, error) {
	var result domain.PlayerAutoConfig
	err := ColPlayerAutoConfig().FindOne(ctx, filterID(userID)).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func ColPlayerAutoConfig() *mongo.Collection {
	return Client.Database(DbName).Collection("playerAutoConfig")
}
