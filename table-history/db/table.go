package db

import (
	"context"
	"github.com/glossd/pokergloss/table-history/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

const collectionName = "tables"

func InsertManyEvents(events []*domain.Event) error {
	if len(events) == 0 {
		return nil
	}
	adaptedEvents := make([]interface{}, 0, len(events))
	for _, table := range events {
		adaptedEvents = append(adaptedEvents, table)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := ColTables().InsertMany(ctx, adaptedEvents)
	if err != nil {
		log.Errorf("Failed to save events: %s", err)
	}
	return err
}

func FindAll() ([]*domain.Event, error) {
	ctx := context.Background()
	var events []*domain.Event

	cur, err := ColTables().Find(ctx, bson.D{})
	if err != nil {
		log.Errorf("Filter events failed: %s", err)
		return nil, err
	}

	for cur.Next(ctx) {
		var e domain.Event
		err := cur.Decode(&e)
		if err != nil {
			return events, err
		}
		events = append(events, &e)
	}

	if err := cur.Err(); err != nil {
		log.Errorf("Filter events failed: %s", err)
		return events, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(events) == 0 {
		return []*domain.Event{}, nil
	}

	return events, nil
}

func DeleteAllBefore(t time.Time) {
	ColTables().DeleteMany(context.Background(), bson.M{"createdat": bson.M{"$lte": bsonx.DateTime(t.UnixNano() / 1e6)}})
}

func ColTables() *mongo.Collection {
	return Client.Database(DbName).Collection(collectionName)
}
