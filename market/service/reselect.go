package service

import (
	"context"
	"github.com/glossd/pokergloss/market/db"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func Reselect(ctx context.Context, userID string) error {
	inv, err := db.FindInventory(ctx, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Errorf("Reselect of non-existing inventory, userID=%s", userID)
			return nil
		}
		return err
	}
	isReselected := inv.Reselect()
	if isReselected {
		err := db.UpsertInventory(ctx, inv)
		if err != nil {
			return err
		}
	}
	return nil
}
