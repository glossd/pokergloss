package mqpub

import (
	"github.com/glossd/memmq"
	"github.com/glossd/pokergloss/market/web/mq"
	log "github.com/sirupsen/logrus"
)

func PublishReselect(userID string) {
	err := memmq.Publish(mq.ReselectTopicID, &mq.ReselectEvent{UserID: userID})
	if err != nil {
		log.Errorf("Failed to publish reselect: %s", err)
	}
}
