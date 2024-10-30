package service

import (
	"context"
	"fmt"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/market/db"
	"github.com/glossd/pokergloss/market/domain"
	"github.com/glossd/pokergloss/market/web/clients/bank"
	"github.com/glossd/pokergloss/market/web/rest/model"
	log "github.com/sirupsen/logrus"
)

func BuyItem(ctx context.Context, iden authid.Identity, params model.BuyItemParams) error {
	command, err := domain.NewPurchaseItemCommand(iden.UserId, domain.ItemID(params.ItemID), params.Units, params.TimeFrame)
	if err != nil {
		return err
	}
	description := fmt.Sprintf("Bought item %s", command.GetItem().Name)
	if command.IsForCoins() {
		err = bank.WithdrawCoins(ctx, command.TotalPrice(), iden.UserId, description)
	} else {
		err = bank.Withdraw(ctx, command.TotalPrice(), iden.UserId, description)
	}
	if err != nil {
		return err
	}

	inventory, err := findOrBuildInventory(ctx, iden.UserId)
	if err != nil {
		return err
	}
	inventory.AddItem(command)

	err = dbBuy(ctx, command, inventory)
	if err != nil {
		refund(command, iden)
		return err
	}
	return nil
}

// TODO Use transaction
func dbBuy(ctx context.Context, command *domain.PurchaseItemCommand, inven *domain.Inventory) error {
	err := db.UpsertInventory(ctx, inven)
	if err != nil {
		log.Errorf("Failed user upsert, commandId=%s", command.ID.Hex())
		return err
	}

	err = db.InsertPurchaseCommand(ctx, command)
	if err != nil {
		log.Errorf("AddItem: failed to insert buy command: %+v", command)
	}
	return nil
}

func refund(command *domain.PurchaseItemCommand, iden authid.Identity) {
	description := fmt.Sprintf("Failed to buy item %s", command.GetItem().Name)
	if command.IsForCoins() {
		bank.DepositCoins(command.TotalPrice(), iden.UserId, description)
	} else {
		bank.Deposit(command.TotalPrice(), iden.UserId, description)
	}
}
