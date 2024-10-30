package balance

import (
	"context"
	"errors"
	"github.com/glossd/pokergloss/bank/db"
	"github.com/glossd/pokergloss/bank/web/mq"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Update(ctx context.Context, userId string, operationId primitive.ObjectID) error {
	op, err := db.FindOperation(ctx, operationId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Errorf("balance.Update: operation not found, opId=%s", operationId.Hex())
			return nil
		}
		log.Errorf("balance.Update: failed to find operation, opId=%s: %s", operationId.Hex(), err)
		return err
	}

	balance, err := db.FindBalance(ctx, userId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Errorf("User balance not found on update: userId=%s", userId)
			return nil
		} else {
			log.Errorf("balance.Update: failed to find balance, userID=%s: %s", userId, err)
			return err
		}
	}

	if balance.LastOperationID == operationId {
		log.Warnf("Tried to update balance with the same last operation id: %s", operationId.Hex())
		return nil
	}

	balance.HandleOperation(op)

	err = db.UpdateBalance(ctx, balance)
	if err != nil {
		log.Errorf("balance.Update: couldn't update balance: %s", err)
		return err
	}

	mq.PublishWsBalanceUpdate(balance)
	return nil
}
