package services

import (
	"context"
	"errors"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/bank/db"
	"github.com/glossd/pokergloss/bank/db/paging"
	"github.com/glossd/pokergloss/bank/domain"
	"github.com/glossd/pokergloss/bank/web/mq/mqpub"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func Deposit(ctx context.Context, operation *domain.Operation) error {
	if operation.Chips <= 0 {
		log.Errorf("Zero-Negative deposit %v", operation)
		return nil
	}
	return doOperation(ctx, operation)
}

func Withdraw(ctx context.Context, operation *domain.Operation) error {
	if operation.Chips <= 0 {
		log.Errorf("Zero-Negative withdraw %v", operation)
		return nil
	}

	b, err := db.FindBalance(ctx, operation.UserID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.ErrNotEnoughChips
		}
		return err
	}

	if !b.IsEnoughChips(operation.Chips) {
		return domain.ErrNotEnoughChips
	}
	return doOperation(ctx, operation)
}

func DepositCoins(ctx context.Context, operation *domain.Operation) error {
	if operation.Coins <= 0 {
		log.Errorf("Zero-Negative withdraw coins %v", operation)
		return nil
	}

	return doOperation(ctx, operation)
}

func WithdrawCoins(ctx context.Context, operation *domain.Operation) error {
	if operation.Coins <= 0 {
		log.Errorf("Zero-Negative withdraw coins %v", operation)
		return nil
	}

	b, err := db.FindBalance(ctx, operation.UserID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.ErrNotEnoughCoins
		}
		return err
	}

	if !b.IsEnoughCoins(operation.Coins) {
		return domain.ErrNotEnoughCoins
	}

	return doOperation(ctx, operation)
}

func doOperation(ctx context.Context, operation *domain.Operation) error {
	err := db.InsertOperation(ctx, operation)
	if err != nil {
		return err
	}

	mqpub.PublishBalanceUpdate(operation.UserID, operation.ID)
	return nil
}

func ListOperations(ctx context.Context, iden authid.Identity, page paging.Page) ([]*domain.Operation, error) {
	return db.FindOperationsByUserIdReverse(ctx, iden.UserId, page)
}
