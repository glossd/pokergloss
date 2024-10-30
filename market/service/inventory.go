package service

import (
	"context"
	"github.com/glossd/pokergloss/market/db"
	"github.com/glossd/pokergloss/market/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func findOrBuildInventory(ctx context.Context, userID string) (*domain.Inventory, error) {
	user, err := db.FindInventory(ctx, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			user = domain.NewInventory(userID)
		} else {
			log.Errorf("db.FindInventory failed: %s", err)
			return nil, err
		}
	}
	return user, nil
}
