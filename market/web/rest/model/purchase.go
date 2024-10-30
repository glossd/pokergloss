package model

import "github.com/glossd/pokergloss/market/domain"

type Purchase struct {
	ID        string           `json:"id"`
	ItemID    domain.ItemID    `json:"itemId"`
	Units     int64            `json:"units"`
	UnitPrice int64            `json:"unitPrice"`
	TimeFrame domain.TimeFrame `json:"timeFrame"`
	CreatedAt int64            `json:"createdAt"`
}

func ToPurchases(bips []*domain.PurchaseItemCommand) []*Purchase {
	result := make([]*Purchase, 0, len(bips))
	for _, bip := range bips {
		result = append(result, ToPurchase(bip))
	}
	return result
}

func ToPurchase(bip *domain.PurchaseItemCommand) *Purchase {
	return &Purchase{
		ID:        bip.ID.Hex(),
		ItemID:    bip.ItemID,
		Units:     bip.Units,
		UnitPrice: bip.UnitPrice,
		TimeFrame: bip.TimeFrame,
		CreatedAt: bip.CreatedAt,
	}
}
