package db

import (
	"context"
	"github.com/glossd/pokergloss/gomq/mqtable"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func InsertGameEnd(ctx context.Context, gameEnd *mqtable.GameEnd) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := ColHistoryGameEnd().InsertOne(ctx, gameEnd)
	if err != nil {
		log.Errorf("Failed insert history game-end: %s", err)
		return err
	}
	return nil
}

func InsertTournamentEnd(ctx context.Context, end *mqtable.TournamentEnd) error {
	_, err := ColHistoryTournamentEnd().InsertOne(ctx, end)
	if err != nil {
		log.Errorf("Failed insert history tournament-end: %s", err)
		return err
	}
	return nil
}

func ForeachGameEnd(apply func(end *mqtable.GameEnd)) error {
	var limit int64 = 1000000
	ctx := context.Background()
	cur, err := ColHistoryGameEnd().Find(ctx, bson.D{}, &options.FindOptions{Limit: &limit})
	if err != nil {
		log.Errorf("Filter tables failed: %s", err)
		return err
	}

	for cur.Next(ctx) {
		var t mqtable.GameEnd
		err := cur.Decode(&t)
		if err != nil {
			log.Errorf("Failed to decode gameEnd: %s", err)
			continue
		}
		apply(&t)
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	return nil
}

func ColHistoryGameEnd() *mongo.Collection {
	return Client.Database(DbName).Collection("historyGameEnd")
}

func ColHistoryTournamentEnd() *mongo.Collection {
	return Client.Database(DbName).Collection("historyTournamentEnd")
}
