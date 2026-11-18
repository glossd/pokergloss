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

// this topic was added to cryptobank after this service was archived.
func SubscribeToMultiDeposits() {
	err := mqbank.PullMultiDeposits("cryptobank-multi-deposits", func(ctx context.Context, r *mqbank.MultiDeposit) error {
		for _, deposit := range r.Deposits {
			// in cryptobank used mongo transactions, all or none
			log.Tracef("Received deposit: %v", r)
			op, err := domain.NewDeposit(domain.Reason(r.Service), deposit.Amount, deposit.UserId, deposit.Description)
			if err != nil {
				log.Errorf("Failed to init deposit operation %+v: %s", deposit, err)
				continue
			}
			err = services.Deposit(ctx, op)
			if err != nil {
				log.Errorf("Failed to save deposit operation %+v: %s", op, err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Couldn't subscribe with pubsub client: %s", err)
	}
}
