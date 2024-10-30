package mqbank

import (
	"context"
	"github.com/glossd/memmq"
	gomq "github.com/glossd/pokergloss/gomq"
	log "github.com/sirupsen/logrus"
)

const BalanceTopicID = "pg.bank.balance-update"

func PublishBalanceUpdate(msg *BalanceUpdate) error {
	return memmq.Publish(BalanceTopicID, msg)
}

func SubscribeForBalanceUpdates(subID string, receiver func(ctx context.Context, msg *BalanceUpdate) error) error {
	return memmq.Subscribe(BalanceTopicID, func(msg interface{}) bool {
		v, ok := msg.(*BalanceUpdate)
		if !ok {
			log.Errorf("Couldn't unmarshal balance update message: %T", v)
			return true
		}
		err := receiver(context.Background(), v)
		return gomq.IsAckableError(err)
	})
}
