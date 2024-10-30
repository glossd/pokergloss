package db

import (
	"context"
	"github.com/glossd/pokergloss/achievement/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindAchievementStore(ctx context.Context, userID string) (*domain.AchievementStore, error) {
	var e domain.AchievementStore
	err := ColAchievementStore().FindOne(ctx, bson.M{"_id": userID}).Decode(&e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func FindAchievementStoreNoCtx(userID string) (*domain.AchievementStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return FindAchievementStore(ctx, userID)
}

func UpsertAchievementStore(ctx context.Context, handCounts *domain.AchievementStore) error {
	_, err := ColAchievementStore().ReplaceOne(ctx, idFilter(handCounts.UserID), handCounts, &options.ReplaceOptions{Upsert: &True})
	if err != nil {
		return err
	}
	return nil
}

func ColAchievementStore() *mongo.Collection {
	return Client.Database(DbName).Collection("achievementStore")
}
