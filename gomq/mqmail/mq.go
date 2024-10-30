package mqmail

import (
	"context"
	"github.com/glossd/memmq"
	"github.com/glossd/pokergloss/gomq"
	log "github.com/sirupsen/logrus"
)

const TopicID = "pg.mail.send"

func Publish(msg *Email) error {
	return memmq.Publish(TopicID, msg)
}

func Subscribe(subID string, receiver func(ctx context.Context, msg *Email) error) error {
	return memmq.Subscribe(TopicID, func(msg interface{}) bool {
		v, ok := msg.(*Email)
		if !ok {
			log.Errorf("memmq: expected *Email, got: %T", v)
			return true
		}
		err := receiver(context.Background(), v)
		return gomq.IsAckableError(err)
	})
}
