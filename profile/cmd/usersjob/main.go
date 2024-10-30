package main

import (
	"context"
	"fmt"
	"github.com/glossd/pokergloss/profile/conf"
	"github.com/glossd/pokergloss/profile/db"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func main() {
	_, _, err := db.Init()
	if err != nil {
		log.Fatalf("db.Init: %s", err)
	}
	ctx := context.Background()
	var nextPageToken string
	usersIter := conf.AuthClient.Users(ctx, nextPageToken)
	for {
		user, err := usersIter.Next()
		if err != nil {
			log.Fatalf("Failed fetch users: %s", err)
		}
		fmt.Println("Handling user " + user.UID)

		_, err = db.Client.Database("firebase").Collection("users").InsertOne(ctx, user)
		if err != nil {
			log.Fatal(err)
		}

		err = conf.AuthClient.DeleteUser(ctx, user.UID)
		if err != nil {
			log.Fatal(err)
		}
	}
}

type Balance struct {
	UserID   string `bson:"_id"`
	Username string
	Picture  string
}

func BalanceCol() *mongo.Collection {
	return db.Client.Database("bank").Collection("balances")
}

func UpdateProfileInfo(b *Balance) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := BalanceCol().UpdateOne(ctx, bson.M{"_id": b.UserID}, bson.M{"$set": bson.M{"username": b.Username, "picture": b.Picture}})
	return err
}
