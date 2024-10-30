package service

import (
	"context"
	"github.com/glossd/pokergloss/market/db"
	"github.com/glossd/pokergloss/market/domain"
)

func ListBuyCommands(ctx context.Context, userID string) ([]*domain.PurchaseItemCommand, error) {
	return db.FindPurchaseCommands(ctx, userID)
}
