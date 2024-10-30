package mqws

import (
	"context"
	"github.com/glossd/memmq"
	log "github.com/sirupsen/logrus"
)

const TableTopicID = "pg.ws.table"

func PublishTableMsg(msg *TableMessage) error {
	return memmq.Publish(TableTopicID, msg)
}

func SubscribeTableMsg(subID string, receiver func(ctx context.Context, msg *TableMessage) error) error {
	return memmq.Subscribe(TableTopicID, func(msg interface{}) bool {
		v, ok := msg.(*TableMessage)
		if !ok {
			log.Errorf("memmq: expected *TableMessage, got: %T", v)
			return true
		}
		err := receiver(context.Background(), v)
		return err == nil
	})
}
