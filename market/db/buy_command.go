package db

import (
	"context"
	"github.com/glossd/pokergloss/market/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertPurchaseCommand(ctx context.Context, cmd *domain.PurchaseItemCommand) error {
	_, err := PurchaseCommandsCol().InsertOne(ctx, cmd)
	if err != nil {
		log.Errorf("Failed to insert buy command: %s", err)
		return err
	}
	return nil
}

func FindPurchaseCommands(ctx context.Context, userID string) ([]*domain.PurchaseItemCommand, error) {
	var commands []*domain.PurchaseItemCommand

	var limit int64 = 10
	cur, err := PurchaseCommandsCol().Find(ctx, bson.M{"userid": userID}, &options.FindOptions{Limit: &limit})
	if err != nil {
		log.Errorf("FindPurchaseCommands: failed to find, userID=%s: %s", userID, err)
		return nil, err
	}

	for cur.Next(ctx) {
		var t domain.PurchaseItemCommand
		err := cur.Decode(&t)
		if err != nil {
			log.Errorf("FindPurchaseCommands: failed to decode, userID=%s: %s", userID, err)
			return commands, err
		}
		commands = append(commands, &t)
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(commands) == 0 {
		return []*domain.PurchaseItemCommand{}, nil
	}

	return commands, nil
}

func PurchaseCommandsCol() *mongo.Collection {
	return Client.Database(DbName).Collection("buyCommands")
}
