package mq

import (
	"github.com/glossd/pokergloss/bank/domain"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/gomq/mqws"
	log "github.com/sirupsen/logrus"
)

func PublishWsBalanceUpdate(balance *domain.Balance) {
	if conf.IsE2E() {
		return
	}
	msg := &mqws.Message{
		ToUserIds: []string{balance.UserID},
		Events:    []*mqws.Event{{Type: "balanceUpdate", Payload: gomq.M{"balance": balance.Chips, "chips": balance.Chips, "coins": balance.Coins}.JSON()}},
	}
	err := mqws.PublishNews(msg)
	if err != nil {
		log.Errorf("Couldn't publish to ws balance update: %v", err)
	} else {
		log.Debugf("Sent balance update: %v", msg)
	}
}
