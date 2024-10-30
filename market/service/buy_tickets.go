package service

import (
	"context"
	"fmt"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/gomq/mqsurvival"
	"github.com/glossd/pokergloss/market/web/clients/bank"
)

const ticketsNum = 5
const ticketsPrice = 5

func BuyTickets(ctx context.Context, iden authid.Identity) error {
	err := bank.WithdrawCoins(ctx, ticketsPrice, iden.UserId, fmt.Sprintf("Buying %d tickets", ticketsNum))
	if err != nil {
		return err
	}
	err = mqsurvival.Publish(&mqsurvival.TicketGift{ToUserId: iden.UserId, Tickets: ticketsNum})
	if err != nil {
		bank.DepositCoins(ticketsPrice, iden.UserId, fmt.Sprintf("Failed to buy %d tickets", ticketsNum))
		return err
	}
	return nil
}
