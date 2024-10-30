package db

import (
	"context"
	"github.com/glossd/pokergloss/achievement/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func FindStatisticsNoCtx(userID string) (*domain.Statistics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return FindStatistics(ctx, userID)
}

func FindStatistics(ctx context.Context, userID string) (*domain.Statistics, error) {
	var s domain.Statistics
	err := ColStatistics().FindOne(ctx, idFilter(userID)).Decode(&s)
	if err != nil {
		log.Errorf("db.FindStatistics failed: %s", err)
		return nil, err
	}
	return &s, nil
}

func UpsertStatistics(s *domain.Statistics) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := ColStatistics().ReplaceOne(ctx, idFilter(s.UserID), s, &options.ReplaceOptions{Upsert: &True})
	if err != nil {
		return err
	}
	return nil
}

func ColStatistics() *mongo.Collection {
	return Client.Database(DbName).Collection("statistics")
}
