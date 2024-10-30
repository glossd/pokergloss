package db

import (
	"context"
	"github.com/glossd/pokergloss/table-chat/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertEmojiMessage(ctx context.Context, msg *domain.EmojiMessage) error {
	_, err := EmojiMsgCol().InsertOne(ctx, msg)
	if err != nil {
		log.Errorf("Failed to insert emoji message: %s", err)
		return err
	}
	return nil
}

func EmojiMsgCol() *mongo.Collection {
	return Client.Database(DbName).Collection("emojis")
}
