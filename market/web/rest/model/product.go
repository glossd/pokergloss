package model

import "github.com/glossd/pokergloss/market/domain"

type Product struct {
	ID       string          `json:"id"`
	SaleType domain.SaleType `json:"saleType" enums:"chips,coins,notForSale"`
	Price    int64
}
