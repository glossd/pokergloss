package mq

import (
	"context"
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/table-history/db"
	"github.com/glossd/pokergloss/table-history/domain"
	log "github.com/sirupsen/logrus"
)

func Subscribe() {
	err := mqws.SubscribeTableMsg("table-history-table-msg", func(ctx context.Context, msg *mqws.TableMessage) error {
		log.Tracef("Got message: %v+", msg)
		return db.InsertManyEvents(domain.ToEvents(msg))
	})
	if err != nil {
		log.Panicf("Failed to subscribe: %s", err)
	}
}
