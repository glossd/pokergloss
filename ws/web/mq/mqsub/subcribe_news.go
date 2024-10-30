package mqsub

import (
	"context"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/ws/web/mq"
	log "github.com/sirupsen/logrus"
)

const NewsSubID = "ws-service"

func SubscribeNews() {
	if !conf.IsE2E() {
		err := mqws.SubscribeNews(NewsSubID, func(ctx context.Context, msg *mqws.Message) error {
			mq.SendMessageToWs(msg)
			return nil
		})
		if err != nil {
			log.Fatalf("mqsub.SubscribeNews failed: %v", err)
		}
		return
	}
}
