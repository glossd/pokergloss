package db

import (
	"context"
	"github.com/glossd/pokergloss/messenger/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindChat(ctx context.Context, id primitive.ObjectID) (*domain.Chat, error) {
	var c domain.Chat
	err := ChatCol().FindOne(ctx, filterID(id)).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func InsertChat(ctx context.Context, chat *domain.Chat) error {
	_, err := ChatCol().InsertOne(ctx, chat)
	if err != nil {
		log.Errorf("Failed to insert chat: %s", err)
		return err
	}
	return nil
}

func ChatCol() *mongo.Collection {
	return Client.Database(DbName).Collection("chats")
}
