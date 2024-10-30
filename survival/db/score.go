package db

import (
	"context"
	"github.com/glossd/pokergloss/survival/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindScoreOrDefault(ctx context.Context, userID string) (*domain.Score, error) {
	var s domain.Score
	err := ScoreCol().FindOne(ctx, filterID(userID)).Decode(&s)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &domain.Score{UserID: userID}, nil
		}
		return nil, err
	}
	return &s, nil
}

func UpsertScore(ctx context.Context, userID string, level int) error {
	_, err := ScoreCol().UpdateOne(ctx, filterID(userID), bson.M{"$max": bson.M{"level": level}}, &options.UpdateOptions{Upsert: &True})
	if err != nil {
		log.Errorf("Failed to upsert score uid=%s, lvl=%d: %s", userID, level, err)
		return err
	}
	return nil
}

func ScoreCol() *mongo.Collection {
	return Client.Database(DbName).Collection("score")
}
