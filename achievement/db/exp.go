package db

import (
	"context"
	"github.com/glossd/pokergloss/achievement/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const DefaultTimeout = time.Second

var True = true

func FindExP(ctx context.Context, userID string) (*domain.ExP, error) {
	var e domain.ExP
	err := ColExP().FindOne(ctx, bson.M{"_id": userID}).Decode(&e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func FindExpNoCtx(userID string) (*domain.ExP, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return FindExP(ctx, userID)
}

func UpsertExpNoCtx(exp *domain.ExP) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	_, err := ColExP().ReplaceOne(ctx, bson.M{"_id": exp.UserID}, exp, &options.ReplaceOptions{Upsert: &True})
	if err != nil {
		log.Errorf("Failed to update ExP=%+v : %s", exp, err)
		return err
	}
	return nil
}

func ColExP() *mongo.Collection {
	return Client.Database(DbName).Collection("points")
}
