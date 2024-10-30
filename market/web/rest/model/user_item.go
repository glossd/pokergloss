package model

import "github.com/glossd/pokergloss/market/domain"

type UserItem struct {
	ItemID domain.ItemID `json:"itemId"`
	// Milliseconds
	ExpiresAt *int64 `json:"expiresAt,omitempty"`
	Selected  bool   `json:"selected"`
}

func ToUserItems(items []*domain.UserItem) []UserItem {
	result := make([]UserItem, 0, len(items))
	for _, item := range items {
		var expiresAt *int64
		if !item.Forever {
			tmp := item.ExpiresAt * 1000
			expiresAt = &tmp
		}
		result = append(result, UserItem{
			ItemID:    item.ItemID,
			ExpiresAt: expiresAt,
			Selected:  item.IsSelected(),
		})
	}
	return result
}
