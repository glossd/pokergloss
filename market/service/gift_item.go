package service

import (
	"context"
	"github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/market/db"
	"github.com/glossd/pokergloss/market/domain"
	log "github.com/sirupsen/logrus"
)

func GiftItem(ctx context.Context, toUserID string, itemID domain.ItemID, units int64, tf domain.TimeFrame) error {
	cmd, err := domain.NewGiftItemCommand(toUserID, itemID, units, tf)
	if err != nil {
		log.Errorf("domain.NewGiftItemCommand error in subscribe for gifts: %s", err)
		return gomq.WrapInAckableError(err)
	}

	inventory, err := findOrBuildInventory(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	inventory.AddItem(cmd)

	err = db.UpsertInventory(ctx, inventory)
	if err != nil {
		log.Errorf("Failed to upsert inventory, userID=%s", cmd.UserID)
		return err
	}
	return nil
}
