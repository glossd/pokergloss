package mqsub

import (
	"context"
	"github.com/glossd/pokergloss/bank/services/balance"
	"github.com/glossd/pokergloss/gomq/mqbank"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SubscribeForBalanceUpdates() {
	err := mqbank.SubscribeForBalanceUpdates("bank-service-balance-updates", func(ctx context.Context, r *mqbank.BalanceUpdate) error {
		id := r.GetOperationId()
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Errorf("Failed to parse operationId: %s", err)
			return nil
		}
		return balance.Update(ctx, r.GetUserId(), oid)
	})
	if err != nil {
		log.Fatalf("Couldn't subscribe with pubsub client: %s", err)
	}
}
