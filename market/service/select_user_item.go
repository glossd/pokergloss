package service

import (
	"context"
	"github.com/glossd/pokergloss/market/db"
	log "github.com/sirupsen/logrus"
)

func SelectUserItem(ctx context.Context, userID string, itemID string) error {
	if userID == "" {
		return E("userID can't be empty")
	}
	if itemID == "" {
		return E("itemID can't be empty")
	}

	user, err := db.FindInventory(ctx, userID)
	if err != nil {
		log.Errorf("SelectUserItem: failed to get user %s: %s", userID, err)
		return err
	}

	err = user.SelectItem(itemID)
	if err != nil {
		return err
	}

	err = db.UpsertInventory(ctx, user)
	if err != nil {
		log.Errorf("SelectUserItem: failed to upsert: %s", err)
		return err
	}

	return nil
}
