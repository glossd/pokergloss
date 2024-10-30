package mqmarket

import (
	"context"
	"github.com/glossd/memmq"
	"github.com/glossd/pokergloss/gomq"
	log "github.com/sirupsen/logrus"
)

const TopicID = "pg.market.gift"

func PublishGift(msg *Gift) error {
	return memmq.Publish(TopicID, msg)
}

func SubscribeForGifts(subID string, receiver func(ctx context.Context, msg *Gift) error) error {
	return memmq.Subscribe(TopicID, func(msg interface{}) bool {
		v, ok := msg.(*Gift)
		if !ok {
			log.Errorf("memmq: expected *Gift, got: %T", v)
			return true
		}
		err := receiver(context.Background(), v)
		return gomq.IsAckableError(err)
	})
}
