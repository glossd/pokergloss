package service

import (
	"context"
	"github.com/glossd/pokergloss/market/domain"
)

func ListUserItems(ctx context.Context, userID string) ([]*domain.UserItem, error) {
	if userID == "" {
		return nil, E("userID can't be empty")
	}

	inv, err := findOrBuildInventory(ctx, userID)
	if err != nil {
		return nil, err
	}

	return inv.AvailableItems(), nil
}
