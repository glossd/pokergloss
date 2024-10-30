package service

import (
	"context"
	"fmt"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/market/domain"
	"github.com/glossd/pokergloss/market/web/rest/model"
)

var ErrNoSuchProduct = E("no such product")

var ProductList = []model.Product{{ID: fmt.Sprintf("%dtickets", ticketsNum), SaleType: domain.ForCoins, Price: ticketsPrice}}

func BuyProduct(ctx context.Context, iden authid.Identity, id string) error {
	switch id {
	case fmt.Sprintf("%dtickets", ticketsNum):
		return BuyTickets(ctx, iden)
	}
	return ErrNoSuchProduct
}
