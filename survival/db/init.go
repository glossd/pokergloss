package db

import (
	"context"
	"errors"
	conf "github.com/glossd/pokergloss/goconf"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var ErrVersionNotMatch = errors.New("version doesn't match")

var Client *mongo.Client

var True = true

func Init() (context.Context, *mongo.Client, error) {
	return InitWithURI(conf.GetDbURI(DbName))
}

func InitWithURI(uri string) (context.Context, *mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return ctx, nil, err
	}
	log.Infof("Connected to MongoDB on %s", conf.Props.DB.Host)

	CreateCardsCol(ctx)

	return ctx, Client, nil
}

func CreateCardsCol(ctx context.Context) {
	names, err := Client.Database(DbName).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		log.Panicf("Failed to list collections: %s", err)
	}
	for _, name := range names {
		if name == cardsCol {
			return
		}
	}

	err = Client.Database(DbName).CreateCollection(ctx, cardsCol, &options.CreateCollectionOptions{
		Validator: bson.M{"$jsonSchema": bson.M{
			"bsonType": "object",
			"properties": bson.M{
				"_id":     bson.M{"bsonType": "string"},
				"tickets": bson.M{"bsonType": "int", "minimum": 0},
			},
		}},
	})
	if err != nil {
		log.Panicf("Failed to create cards collection: %s", err)
	}
}

func filterID(id interface{}) bson.M {
	return bson.M{"_id": id}
}
