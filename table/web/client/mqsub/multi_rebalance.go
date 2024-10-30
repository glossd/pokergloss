package mqsub

import (
	"context"
	"github.com/glossd/memmq"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	"github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/table/services/multi"
	"github.com/glossd/pokergloss/table/web/client/mq"
	log "github.com/sirupsen/logrus"
)

func SubscribeForMultiRebalance() {
	if conf.IsE2E() {
		for event := range mq.TestMultiPlayersMovedQueue {
			if mq.IsMultiPlayersMovedEnabledTest {
				ctx := context.Background()
				err := multi.CountDownRebalance(ctx, event)
				if err != nil {
					log.Panic(err)
				}
			}
		}
		return
	}
	err := memmq.Subscribe(mq.MultiPlayersMovedTopicID, func(message interface{}) bool {
		event, ok := message.(*mq.MultiPlayersMovedEvent)
		if !ok {
			log.Errorf("Failed to parse start multi event, got=%T", event)
			return true
		}
		<-timeutil.AfterTimeAt(event.RebalanceAt)
		err := multi.CountDownRebalance(context.Background(), event)
		return gomq.IsAckableError(err)
	})
	if err != nil {
		log.Panicf("Failed to init timeout subscriber: %s", err)
	}
}
