package mqpub

import (
	"github.com/glossd/memmq"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	"github.com/glossd/pokergloss/table/web/client/mq"
	log "github.com/sirupsen/logrus"
)

var timeoutPublisher *gomq.Publisher

func InitTimeoutPublisher() {
	if conf.IsProd() {
		var err error
		timeoutPublisher, err = gomq.InitPublisher()
		if err != nil {
			log.Panicf("Failed to init timeout timeoutPublisher: %s", err)
		}
	}
}

func PublishTimeoutEvent(event *timeout.Event) {
	mq.SetCacheGameFlow(event.Key.TableID, event.Key.Version)
	if conf.IsE2E() {
		if mq.IsTimeoutTestMQEnabled {
			mq.TimeoutTestMQ <- event
		}
		return
	}

	if conf.IsProd() {
		if timeoutPublisher == nil {
			log.Errorf("Timeout timeoutPublisher is nil")
			return
		}

		err := timeoutPublisher.PublishJSON(mq.TimeoutTopicID, event)
		if err != nil {
			log.Errorf("Failed to send timeout event: %s", err)
		}

	} else {
		err := memmq.Publish(mq.TimeoutTopicID, event)
		if err != nil {
			log.Errorf("Failed to send timeout event: %s", err)
		}
	}
}
