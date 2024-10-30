package db

import (
	"context"
	"github.com/glossd/pokergloss/messenger/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MaxLimit = 20

func FindMessage(ctx context.Context, msgID primitive.ObjectID) (*domain.Message, error) {
	var msg domain.Message
	err := MessageCol().FindOne(ctx, filterID(msgID)).Decode(&msg)
	if err != nil {
		return nil, err
	}
	return &msg, err
}

func FindMessagesByChatID(ctx context.Context, chatID primitive.ObjectID, lastOID primitive.ObjectID, limit int64) ([]*domain.Message, error) {
	filter := bson.D{{"chatid", chatID}}
	if lastOID != primitive.NilObjectID {
		filter = bson.D{{"chatid", chatID}, {"_id", bson.M{"$lt": lastOID}}}
	}
	if limit > MaxLimit || limit < 1 {
		limit = MaxLimit
	}
	cur, err := MessageCol().Find(ctx, filter, &options.FindOptions{Limit: &limit, Sort: bson.M{"_id": -1}})
	if err != nil {
		return nil, err
	}
	var messages []*domain.Message
	for cur.Next(ctx) {
		var m domain.Message
		err := cur.Decode(&m)
		if err != nil {
			return messages, err
		}
		messages = append(messages, &m)
	}

	if err := cur.Err(); err != nil {
		log.Errorf("Find messages failed: %s", err)
		return messages, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(messages) == 0 {
		return []*domain.Message{}, nil
	}

	return messages, nil
}

func FindLastMessageWithChatID(ctx context.Context, chatID primitive.ObjectID) (*domain.Message, error) {
	var one int64 = 1
	cur, err := MessageCol().Find(ctx, bson.M{"chatid": chatID}, &options.FindOptions{Limit: &one, Sort: bson.M{"_id": -1}})
	if err != nil {
		return nil, err
	}
	if cur.Next(ctx) {
		var result domain.Message
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		} else {
			cur.Close(ctx)
			return &result, nil
		}
	}
	return nil, mongo.ErrNoDocuments
}

func InsertMessage(ctx context.Context, msg *domain.Message) error {
	_, err := MessageCol().InsertOne(ctx, msg)
	if err != nil {
		log.Errorf("InsertMessage failed: %s", err)
		return err
	}
	return nil
}

func UpdateMessageStatus(ctx context.Context, msgID primitive.ObjectID, status domain.MessageStatus) error {
	_, err := MessageCol().UpdateOne(ctx, filterID(msgID), bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		log.Errorf("UpdateMessageStatus failed: %s", err)
		return err
	}
	return nil
}

func MessageCol() *mongo.Collection {
	return Client.Database(DbName).Collection("messages")
}
