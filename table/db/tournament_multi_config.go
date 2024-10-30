package db

import (
	"context"
	"errors"
	"github.com/glossd/pokergloss/goconf/timeutil"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const id = "last"

type MultiConfig struct {
	ID                       string `bson:"_id"`
	CreatedHourlyFreerollsAt int64
}

func GetMultiConfig() (*MultiConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	var config MultiConfig
	err := ColLobbyMultiConfig().FindOne(ctx, filterMultiConfigID()).Decode(&config)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	return &config, nil
}

func UpdateLastCreatedHourlyFreerolls(t time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	upsert := true
	_, err := ColLobbyMultiConfig().UpdateOne(ctx, filterMultiConfigID(),
		bson.M{"$set": bson.M{"createdhourlyfreerollsat": timeutil.Time(t)}}, &options.UpdateOptions{Upsert: &upsert})
	if err != nil {
		log.Errorf("Failed update multi config CreatedHourlyFreerollsAt: %s", err)
	}
}

func ColLobbyMultiConfig() *mongo.Collection {
	return Client.Database(DbName).Collection("multiConfig")
}

func filterMultiConfigID() bson.M {
	return bson.M{"_id": id}
}
