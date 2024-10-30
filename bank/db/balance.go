package db

import (
	"context"
	"github.com/glossd/pokergloss/bank/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertBalance(ctx context.Context, balance *domain.Balance) error {
	_, err := BalanceCol().InsertOne(ctx, balance)
	if err != nil {
		return err
	}
	return nil
}

func UpdateBalanceNoCtx(balance *domain.Balance) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return UpdateBalance(ctx, balance)
}

func UpdateBalance(ctx context.Context, balance *domain.Balance) error {
	balance.Version++
	_, err := BalanceCol().ReplaceOne(ctx, bson.D{{"_id", balance.UserID}, {"version", balance.Version - 1}}, balance)
	if err != nil {
		return err
	}
	return nil
}

func UpdateProfileInfo(ctx context.Context, b *domain.Balance) error {
	_, err := BalanceCol().UpdateOne(ctx, bson.M{"_id": b.UserID}, bson.M{"$set": bson.M{"username": b.Username, "picture": b.Picture}})
	return err
}

func FindBalanceNoCtx(userId string) (*domain.Balance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return FindBalance(ctx, userId)
}

func FindBalance(ctx context.Context, userId string) (*domain.Balance, error) {
	one := BalanceCol().FindOne(ctx, filterStrID(userId))
	var b domain.Balance
	err := one.Decode(&b)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func ForEachBalance(ctx context.Context, apply func(balance *domain.Balance)) error {
	cur, err := BalanceCol().Find(ctx, bson.D{})
	if err != nil {
		log.Errorf("ForEachBalance failed: %s", err)
		return err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var b domain.Balance
		err := cur.Decode(&b)
		if err != nil {
			log.Errorf("Failed to decode balance: %s", err)
			return err
		}
		apply(&b)
	}
	return nil
}

func filterStrID(id string) bson.D {
	return bson.D{{"_id", id}}
}

func filterID(id primitive.ObjectID) bson.D {
	return bson.D{{"_id", id}}
}

func BalanceCol() *mongo.Collection {
	return Client.Database(DbName).Collection("balances")
}
