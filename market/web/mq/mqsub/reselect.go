package mqsub

import (
	"context"
	"github.com/glossd/memmq"
	"github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/market/service"
	"github.com/glossd/pokergloss/market/web/mq"
	log "github.com/sirupsen/logrus"
)

func SubscribeForReselect() {
	err := memmq.Subscribe(mq.ReselectTopicID, func(msg interface{}) bool {
		v, ok := msg.(*mq.ReselectEvent)
		if !ok {
			log.Errorf("Failed to parse reselect event: %T", v)
			return true
		}
		err := service.Reselect(context.Background(), v.UserID)
		return gomq.IsAckableError(err)
	})
	if err != nil {
		log.Panicf("Failed to init timeout subscriber: %s", err)
	}
}
