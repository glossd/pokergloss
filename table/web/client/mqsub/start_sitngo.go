package mqsub

import (
	"context"
	"github.com/glossd/memmq"
	"github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/table/services/sitngo"
	"github.com/glossd/pokergloss/table/web/client/mq"
	log "github.com/sirupsen/logrus"
)

func SubscribeForStartSitngo() {
	err := memmq.Subscribe(mq.StartSitngoTopicID, func(msg interface{}) bool {
		v, ok := msg.(*mq.StartTournamentEvent)
		if !ok {
			log.Errorf("Failed to parse start tournament event: %T", v)
			return true
		}
		err := sitngo.Start(context.Background(), v.StartAt)
		return gomq.IsAckableError(err)
	})
	if err != nil {
		log.Panicf("Failed to init timeout subscriber: %s", err)
	}
}
