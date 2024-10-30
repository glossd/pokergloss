package balance

import (
	"context"
	"github.com/glossd/pokergloss/bank/db"
	"github.com/glossd/pokergloss/bank/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

func Get(ctx context.Context, userId string) (*domain.Balance, error) {
	balance, err := db.FindBalance(ctx, userId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &domain.Balance{UserID: userId}, nil
		}
		return nil, err
	}
	return balance, nil
}
