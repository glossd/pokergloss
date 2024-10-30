package db

import (
	"context"
	"github.com/glossd/pokergloss/survival/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func Find(ctx context.Context, userID string) (*domain.Survival, error) {
	var s domain.Survival
	err := SurvivalCol().FindOne(ctx, filterID(userID)).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func Insert(ctx context.Context, s *domain.Survival) error {
	_, err := SurvivalCol().InsertOne(ctx, s)
	if err != nil {
		log.Errorf("Insert survival failed: %s", err)
		return err
	}
	return nil
}

func Update(ctx context.Context, s *domain.Survival) error {
	_, err := SurvivalCol().ReplaceOne(ctx, filterID(s.UserID), s)
	if err != nil {
		log.Errorf("Update survival failed: %s", err)
		return err
	}
	return nil
}

func Delete(ctx context.Context, userID string) error {
	_, err := SurvivalCol().DeleteOne(ctx, filterID(userID))
	if err != nil {
		log.Errorf("Delete survival failed: %s", err)
		return err
	}
	return nil
}

func SurvivalCol() *mongo.Collection {
	return Client.Database(DbName).Collection("survivals")
}
