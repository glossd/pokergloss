package db

import (
	"context"
	"github.com/glossd/pokergloss/table-chat/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "messages"

func InsertMessage(ctx context.Context, msg *domain.Message) {
	_, err := MessageCol().InsertOne(ctx, msg)
	if err != nil {
		log.Errorf("Failed to insert message: %s", err)
	}
}

func MessageCol() *mongo.Collection {
	return Client.Database(DbName).Collection(collectionName)
}
