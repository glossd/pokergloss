package mqsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/glossd/memmq"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	"github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	"github.com/glossd/pokergloss/table/web/client/mq"
	log "github.com/sirupsen/logrus"
)

func SubscribeForTimeouts() {
	if conf.IsE2E() {
		for event := range mq.TimeoutTestMQ {
			if mq.IsTimeoutTestMQEnabled {
				nack := HandleTimeoutEvent(context.Background(), *event)
				if nack {
					mq.TimeoutTestMQ <- event
				}
			}
		}
		return
	}

	if conf.IsProd() {
		err := gomq.Pull("table-service-timeout", mq.TimeoutTopicID, func(ctx context.Context, msg *pubsub.Message) error {
			var event timeout.Event
			err := json.Unmarshal(msg.Data, &event)
			if err != nil {
				log.Errorf("Failed to parse timeout message: %s", err)
				return nil
			}

			nack := HandleTimeoutEvent(ctx, event)
			if nack {
				return fmt.Errorf("handle timeout error")
			} else {
				return nil
			}
		})
		if err != nil {
			log.Panicf("Failed to init timeout subscriber: %s", err)
		}
	} else {
		err := memmq.Subscribe(mq.TimeoutTopicID, func(msg interface{}) bool {
			v, ok := msg.(*timeout.Event)
			if !ok {
				log.Errorf("Failed to parse timeout message: %T", v)
				return true
			}
			nack := HandleTimeoutEvent(context.Background(), *v)
			return !nack
		})
		if err != nil {
			log.Panicf("Failed to init memmq timeout subscriber: %s", err)
		}
	}
}

func HandleTimeoutEvent(ctx context.Context, event timeout.Event) (nack bool) {
	log.Debugf("Got timeout event: %+v", event)
	select {
	case <-ctx.Done():
		log.Errorf("Context done for timeout event=%+v", event)
		return true
	case <-timeutil.AfterTimeAt(event.At):
		switch event.Type {
		case timeout.SeatReservation:
			return actionhandler.DoSeatReservationTimeout(ctx, event.Key)
		case timeout.Decision:
			gfv := mq.GetCacheGameFlow(event.Key.TableID)
			if gfv > event.Key.Version {
				return false
			}
			return actionhandler.DoDecisionTimeout(ctx, event.Key)
		case timeout.StartGame:
			return actionhandler.DoStartGame(ctx, event.Key)
		}
	}
	return true
}
