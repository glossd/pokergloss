package mqsub

import (
	"context"
	"github.com/glossd/memmq"
	"github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/table/services/multi"
	"github.com/glossd/pokergloss/table/web/client/mq"
	log "github.com/sirupsen/logrus"
)

func SubscribeForStartMulti() {
	err := memmq.Subscribe(mq.StartMultiTopicID, func(msg interface{}) bool {
		v, ok := msg.(*mq.StartTournamentEvent)
		if !ok {
			log.Errorf("memmq: expected *GameEnd, got: %T", v)
			return true
		}
		err := multi.StartMultiLobbies(context.Background(), v.StartAt)
		return gomq.IsAckableError(err)
	})
	if err != nil {
		log.Panicf("Failed to init timeout subscriber: %s", err)
	}
}
