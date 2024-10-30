package mqpub

import (
	"github.com/glossd/memmq"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	"github.com/glossd/pokergloss/table/web/client/mq"
	log "github.com/sirupsen/logrus"
)

func PublishTimeoutEvent(event *timeout.Event) {
	if conf.IsE2E() {
		if mq.IsTimeoutTestMQEnabled {
			mq.TimeoutTestMQ <- event
		}
		return
	}

	err := memmq.Publish(mq.TimeoutTopicID, event)
	if err != nil {
		log.Errorf("Failed to send timeout event: %s", err)
	}
}
