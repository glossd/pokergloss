package mqws

import (
	"context"
	"github.com/glossd/memmq"
	log "github.com/sirupsen/logrus"
)

// todo, rename `all` to `news`
const TopicID = "pg.ws.all"

// Deprecated, use PublishNews
func Publish(msg *Message) error {
	return PublishNews(msg)
}

func PublishNews(msg *Message) error {
	return memmq.Publish(TopicID, msg)
}

// Deprecated, use SubscribeNews
func Subscribe(subID string, receiver func(ctx context.Context, msg *Message) error) error {
	return SubscribeNews(subID, receiver)
}

func SubscribeNews(subID string, receiver func(ctx context.Context, msg *Message) error) error {
	return memmq.Subscribe(TopicID, func(msg interface{}) bool {
		v, ok := msg.(*Message)
		if !ok {
			log.Errorf("memmq: expected *GameEnd, got: %T", v)
			return true
		}
		err := receiver(context.Background(), v)
		return err == nil
	})
}
