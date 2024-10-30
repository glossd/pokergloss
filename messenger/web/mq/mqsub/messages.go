package mqsub

import (
	"context"
	"github.com/glossd/pokergloss/gomq/mqmessenger"
	"github.com/glossd/pokergloss/messenger/service"
	log "github.com/sirupsen/logrus"
)

func Subscribe() {
	err := mqmessenger.SubscribeForMessages("messenger-service", func(ctx context.Context, msg *mqmessenger.Message) error {
		return service.InnerSendMessage(ctx, msg.FromUserId, msg.ToUserId, msg.Text)
	})
	if err != nil {
		log.Panicf("Failed to subscribe messenger mq: %s", err)
	}
}
