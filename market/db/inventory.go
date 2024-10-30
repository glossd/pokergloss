package db

import (
	"context"
	"github.com/glossd/pokergloss/market/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindInventoryNoCtx(userID string) (*domain.Inventory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return FindInventory(ctx, userID)
}

func FindInventory(ctx context.Context, userID string) (*domain.Inventory, error) {
	var user domain.Inventory
	err := InventoryCol().FindOne(ctx, filterID(userID)).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpsertInventory(ctx context.Context, inventory *domain.Inventory) error {
	if inventory.Version == 0 {
		inventory.Version++
		_, err := InventoryCol().InsertOne(ctx, inventory)
		if err != nil {
			log.Errorf("Failed to insert inventory: %s", err)
			return err
		}
	} else {
		filter := filterIdAndVersion(inventory.UserID, inventory.Version)
		inventory.Version++
		_, err := InventoryCol().ReplaceOne(ctx, filter, inventory)
		if err != nil {
			log.Errorf("Failed to update inventory: %s", err)
			return err
		}
	}
	return nil
}

func ForeachInventory(apply func(u *domain.Inventory)) error {
	ctx := context.Background()
	cur, err := InventoryCol().Find(ctx, bson.D{})
	if err != nil {
		return err
	}
	for cur.Next(ctx) {
		var u domain.Inventory
		err := cur.Decode(&u)
		if err != nil {
			return err
		}
		apply(&u)
	}
	cur.Close(ctx)
	return nil
}

func InventoryCol() *mongo.Collection {
	return Client.Database(DbName).Collection("users")
}
