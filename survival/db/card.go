package db

import (
	"context"
	"github.com/glossd/pokergloss/survival/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const cardsCol = "cards"

func FindCard(ctx context.Context, userID string) (*domain.Card, error) {
	var c domain.Card
	err := CardCol().FindOne(ctx, filterID(userID)).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func CardIncTicket(ctx context.Context, userID string) error {
	err := updateCard(ctx, userID, 1)
	if err != nil {
		log.Errorf("Failed to increment card tickets: %s", err)
		return err
	}
	return nil
}

func CardIncTwoTickets(ctx context.Context, userID string) error {
	err := updateCard(ctx, userID, 2)
	if err != nil {
		log.Errorf("Failed to increment card tickets: %s", err)
		return err
	}
	return nil
}

func CardDecTicket(ctx context.Context, userID string) error {
	return updateCard(ctx, userID, -1)
}

func GiftTickets(ctx context.Context, userID string, tickets int64) error {
	err := updateCard(ctx, userID, int(tickets))
	if err != nil {
		log.Errorf("Failed to Gift tickets userID=%s: %s", userID, err)
		return err
	}
	return nil
}

func updateCard(ctx context.Context, userID string, add int) error {
	_, err := CardCol().UpdateOne(ctx, filterID(userID), bson.M{"$inc": bson.M{"tickets": add}}, &options.UpdateOptions{Upsert: &True})
	return err
}

func CardCol() *mongo.Collection {
	return Client.Database(DbName).Collection(cardsCol)
}
