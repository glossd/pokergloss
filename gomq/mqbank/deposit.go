package mqbank

import (
	"context"
	"github.com/glossd/memmq"
	log "github.com/sirupsen/logrus"
)

const TopicID = "pg.bank.deposits"
const MultiDepositTopicID = "pocrium.bank.multi-deposits"

func Deposit(request *DepositRequest) error {
	return memmq.Publish(TopicID, request)
}

func PullDeposits(subID string, process func(ctx context.Context, r *DepositRequest) error) error {
	return memmq.Subscribe(TopicID, func(msg interface{}) bool {
		d, ok := msg.(*DepositRequest)
		if !ok {
			log.Errorf("memmq: expected *DepositRequest, got: %T", d)
			return true
		}
		err := process(context.Background(), d)
		return err == nil
	})
}

func PublishMultiDeposit(ctx context.Context, r *MultiDeposit) error {
	return memmq.Publish(MultiDepositTopicID, r)
}

func PublishMultiDepositAsync(r *MultiDeposit) error {
	return memmq.Publish(MultiDepositTopicID, r)
}

func PullMultiDeposits(subID string, process func(ctx context.Context, r *MultiDeposit) error) error {
	return memmq.Subscribe(MultiDepositTopicID, func(msg interface{}) bool {
		d, ok := msg.(*MultiDeposit)
		if !ok {
			log.Errorf("memmq: expected *MultiDeposit, got: %T", d)
			return true
		}
		err := process(context.Background(), d)
		return err == nil
	})
}
