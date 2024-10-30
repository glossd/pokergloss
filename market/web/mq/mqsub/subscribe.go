package mqsub

import (
	"context"
	"github.com/glossd/pokergloss/gomq/mqmarket"
	"github.com/glossd/pokergloss/market/domain"
	"github.com/glossd/pokergloss/market/service"
	log "github.com/sirupsen/logrus"
)

func Subscribe() {
	err := mqmarket.SubscribeForGifts("market-service", func(ctx context.Context, msg *mqmarket.Gift) error {
		return service.GiftItem(ctx, msg.ToUserId, domain.ItemID(msg.ItemId), msg.Units, toTimeFrame(msg.TimeFrame))
	})
	if err != nil {
		log.Panicf("Failed to subscribe for gifts: %s", err)
	}
}

func toTimeFrame(tf mqmarket.Gift_TimeFrame) domain.TimeFrame {
	switch tf {
	case mqmarket.Gift_DAY:
		return domain.Day
	case mqmarket.Gift_WEEK:
		return domain.Week
	case mqmarket.Gift_MONTH:
		return domain.Month
	}
	return domain.Day
}
