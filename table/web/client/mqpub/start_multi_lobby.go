package mqpub

import (
	"github.com/glossd/memmq"
	"github.com/glossd/pokergloss/table/web/client/mq"
	log "github.com/sirupsen/logrus"
)

func PublishStartMulti(startAt int64) {
	event := &mq.StartTournamentEvent{StartAt: startAt}

	err := memmq.Publish(mq.StartMultiTopicID, event)
	if err != nil {
		log.Errorf("Failed to send timeout event: %s", err)
	}

}
