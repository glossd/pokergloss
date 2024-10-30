package mqsurvival

import (
	"context"
	"github.com/glossd/memmq"
	"github.com/glossd/pokergloss/gomq"
	log "github.com/sirupsen/logrus"
)

const TopicID = "pg.survival.ticket-gift"

func Publish(msg *TicketGift) error {
	return memmq.Publish(TopicID, msg)
}

func SubscribeForTicketGifts(subID string, receiver func(ctx context.Context, msg *TicketGift) error) error {
	return memmq.Subscribe(TopicID, func(msg interface{}) bool {
		v, ok := msg.(*TicketGift)
		if !ok {
			log.Errorf("memmq: expected *TicketGift, got: %T", v)
			return true
		}
		err := receiver(context.Background(), v)
		return gomq.IsAckableError(err)
	})
}
