package model

import "github.com/glossd/pokergloss/market/domain"

type Item struct {
	ID           domain.ItemID       `json:"id"`
	SaleType     domain.SaleType     `json:"saleType" enums:"chips,coins,notForSale"`
	PositionType domain.PositionType `json:"positionType" enums:"side,top"`
	PriceList    PriceList           `json:"priceList"`
}

type PriceList struct {
	Day   int64 `json:"day"`
	Week  int64 `json:"week"`
	Month int64 `json:"month"`
}

func ToItems(items []*domain.Item) []*Item {
	result := make([]*Item, 0, len(items))
	for _, product := range items {
		result = append(result, ToItem(product))
	}
	return result
}

func ToItem(p *domain.Item) *Item {
	return &Item{
		ID:           p.ID,
		SaleType:     p.SaleType,
		PositionType: p.PositionType,
		PriceList: PriceList{
			Day:   p.PriceList.Day,
			Week:  p.PriceList.Week,
			Month: p.PriceList.Month,
		},
	}
}
