package mqpub

import (
	"context"
	"github.com/glossd/pokergloss/bank/services/balance"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq/mqbank"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PublishBalanceUpdate(userId string, operationId primitive.ObjectID) {
	if conf.IsE2E() {
		err := balance.Update(context.Background(), userId, operationId)
		if err != nil {
			log.Fatalf("BalanceUpdate: %s", err)
		}
		return
	}
	err := mqbank.PublishBalanceUpdate(&mqbank.BalanceUpdate{UserId: userId, OperationId: operationId.Hex()})
	if err != nil {
		log.Errorf("Failed to send balance update: %s", err)
	}
}
