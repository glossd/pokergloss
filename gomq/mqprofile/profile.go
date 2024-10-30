package mqprofile

import (
	"context"
	"github.com/glossd/memmq"
	"github.com/glossd/pokergloss/gomq"
	log "github.com/sirupsen/logrus"
)

const TopicID = "pg.profile.updates"
const CreatedTopicID = "pg.profile.created"

func Publish(msg *Profile) error {
	return memmq.Publish(TopicID, msg)
}

func PublishCreated(msg *Profile) error {
	return memmq.Publish(CreatedTopicID, msg)
}

func Subscribe(subID string, receiver func(ctx context.Context, msg *Profile) error) error {
	return memmq.Subscribe(TopicID, func(msg interface{}) bool {
		v, ok := msg.(*Profile)
		if !ok {
			log.Errorf("Couldn't unmarshal profile message: %T", v)
			return true
		}
		err := receiver(context.Background(), v)
		return gomq.IsAckableError(err)
	})
}

func SubscribeForCreated(subID string, receiver func(ctx context.Context, msg *Profile) error) error {
	return memmq.Subscribe(CreatedTopicID, func(msg interface{}) bool {
		v, ok := msg.(*Profile)
		if !ok {
			log.Errorf("memmq: expected *Profile, got: %T", v)
			return true
		}
		err := receiver(context.Background(), v)
		return err == nil
	})
}
