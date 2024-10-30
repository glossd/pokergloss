package mqpub

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq/mqws"
	log "github.com/sirupsen/logrus"
)

func Publish(msg *mqws.TableMessage) {
	if conf.IsProd() {
		err := mqws.PublishTableMsg(msg)
		if err != nil {
			log.Errorf("Failed to publish message to pubsub: %s", err)
		}
		return
	}
}
