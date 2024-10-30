package db

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Used to exist to limit the number of k8s pods
const counterID = "anonymous"
const DbName = "survival"

type Counter struct {
	ID    string `bson:"_id"`
	Count int
}

func FindAnonymousCounter(ctx context.Context) (int, error) {
	var counter Counter
	err := AnonymousCounterCol().FindOne(ctx, filterID(counterID)).Decode(&counter)
	if err == mongo.ErrNoDocuments {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return counter.Count, nil
}

func IncAnonymousCounter(ctx context.Context) error {
	return updateCounter(ctx, 1)
}

func DecAnonymousCounter(ctx context.Context) error {
	return updateCounter(ctx, -1)
}

func updateCounter(ctx context.Context, count int) error {
	_, err := AnonymousCounterCol().UpdateOne(ctx, filterID(counterID), bson.M{"$inc": bson.M{"count": count}}, &options.UpdateOptions{Upsert: &True})
	if err != nil {
		log.Errorf("Failed to increment anonymous counter: %s", err)
		return err
	}
	return nil
}

func AnonymousCounterCol() *mongo.Collection {
	return Client.Database(DbName).Collection("anonymousCounter")
}
