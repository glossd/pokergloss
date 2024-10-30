package model

import (
	"github.com/glossd/pokergloss/market/domain"
)

type BuyItemParams struct {
	ItemID    string           `json:"itemId"`
	Units     int64            `json:"units"`
	TimeFrame domain.TimeFrame `json:"timeFrame" enums:"day,week,month"`
}
