package mqsub

import (
	"context"
	"github.com/glossd/pokergloss/bank/domain"
	"github.com/glossd/pokergloss/bank/services"
	"github.com/glossd/pokergloss/bank/web"
	"github.com/glossd/pokergloss/gomq/mqbank"
	log "github.com/sirupsen/logrus"
)

func SubscribeToDeposits() {
	err := mqbank.PullDeposits("", func(ctx context.Context, r *mqbank.DepositRequest) error {
		log.Tracef("Received deposit: %v", r)
		operation, err := web.ToOperation(r)
		if err != nil {
			return err
		}
		if operation.Type == domain.Deposit {
			return services.Deposit(ctx, operation)
		} else if operation.Type == domain.DepositCoins {
			return services.DepositCoins(ctx, operation)
		}
		log.Errorf("No such deposit operation type: %s", operation.Type)
		return nil
	})
	if err != nil {
		log.Fatalf("Couldn't subscribe with pubsub client: %s", err)
	}
}
