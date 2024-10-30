package mqmessenger

import (
	"context"
	"github.com/glossd/memmq"
	"github.com/glossd/pokergloss/gomq"
	log "github.com/sirupsen/logrus"
)

const TopicID = "pg.messenger.send"

func Publish(msg *Message) error {
	return memmq.Publish(TopicID, msg)
}

func SubscribeForMessages(subID string, receiver func(ctx context.Context, msg *Message) error) error {
	return memmq.Subscribe(TopicID, func(msg interface{}) bool {
		v, ok := msg.(*Message)
		if !ok {
			log.Errorf("memmq: expected *Message, got: %T", v)
			return true
		}
		err := receiver(context.Background(), v)
		return gomq.IsAckableError(err)
	})
}
