package mqpub

import (
	"fmt"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq/mqbank"
	"github.com/glossd/pokergloss/gomq/mqmarket"
	"github.com/glossd/pokergloss/survival/domain"
	log "github.com/sirupsen/logrus"
)

func Deposit(s *domain.Survival) error {
	if conf.IsE2E() {
		return nil
	}
	wof := s.GetWheelOfFortune()
	won := wof.WonSlot()
	dep := &mqbank.DepositRequest{
		Chips:       won.Chips,
		Type:        mqbank.DepositRequest_SURVIVAL,
		Description: fmt.Sprintf("Reached %d level in Survival", s.Level),
		UserId:      s.UserID,
	}
	err := mqbank.Deposit(dep)
	if err != nil {
		log.Errorf("Failed to deposit %+v : %s", dep, err)
		return err
	}
	return nil
}

func GiftItem(s *domain.Survival) error {
	wof := s.GetWheelOfFortune()
	item := wof.WonSlot().Item
	gift := &mqmarket.Gift{
		ItemId:    item.ItemID,
		TimeFrame: mqmarket.Gift_DAY,
		Units:     item.Days,
		ToUserId:  s.UserID,
	}
	err := mqmarket.PublishGift(gift)
	if err != nil {
		log.Errorf("Failed to send market item %+v : %s", gift, err)
		return err
	}
	return nil
}
