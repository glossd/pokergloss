package service

import (
	"context"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/market/domain"
	"github.com/glossd/pokergloss/market/web/mq/mqpub"
)

func GetSelectedItem(ctx context.Context, userID string) (*domain.UserItem, error) {
	if userID == "" {
		return nil, E("userId, can't be null")
	}

	inv, err := findOrBuildInventory(ctx, userID)
	if err != nil {
		return nil, err
	}

	selected, err := inv.GetSelectedItem()
	if err != nil {
		if err == domain.ErrSelectItemExpired {
			if inv.Reselect() { // should be always true
				if conf.IsE2E() {
					_ = Reselect(ctx, userID)
				} else {
					mqpub.PublishReselect(userID)
				}
				return inv.GetSelectedItem()
			}
		}
		return nil, err
	}

	return selected, nil
}
