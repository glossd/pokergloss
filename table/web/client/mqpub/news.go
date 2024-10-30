package mqpub

import (
	conf "github.com/glossd/pokergloss/goconf"
	mqws2 "github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/web/client/mq"
	log "github.com/sirupsen/logrus"
)

func SendNewsToUsers(userIDs []string, event *events.TableEvent) {
	msg := &mqws2.Message{
		ToUserIds: userIDs,
		Events:    []*mqws2.Event{ToEvent(event)},
	}
	publishNews(msg)
}

func publishNews(msg *mqws2.Message) {
	if conf.IsProd() || conf.IsLocalOnly() {
		err := mqws2.PublishNews(msg)
		if err != nil {
			log.Errorf("Failed to publish news message: %v", err)
		} else {
			log.Tracef("Published mqws.Message msg: %v", msg)
		}
	} else if conf.IsE2E() {
		mq.TestNewsMQ <- msg
	}
}
