package db

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var True = true

type Notification struct {
	UserID string `bson:"_id"`
	Web    Web
}

type Web struct {
	Token string
}

func FindNotification(ctx context.Context, userID string) (*Notification, error) {
	var n Notification
	err := NotificationCol().FindOne(ctx, filterID(userID)).Decode(&n)
	if err != nil {
		return nil, err
	}
	return &n, err
}

func UpsertNotification(ctx context.Context, n *Notification) error {
	_, err := NotificationCol().ReplaceOne(ctx, filterID(n.UserID), n, &options.ReplaceOptions{Upsert: &True})
	if err != nil {
		log.Errorf("Failed to upsert notification of userID=%s : %s", n.UserID, err)
		return err
	}
	return nil
}

func NotificationCol() *mongo.Collection {
	return Client.Database(DbName).Collection("notifications")
}
