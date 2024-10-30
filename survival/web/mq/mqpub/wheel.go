package mqpub

import (
	conf "github.com/glossd/pokergloss/goconf"
	gomq "github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/survival/domain"
	log "github.com/sirupsen/logrus"
)

type WoF struct {
	Slots      []Slot `json:"slots"`
	WonSlotIdx int    `json:"wonSlotIdx"`
}

type Slot struct {
	Chips *int64 `json:"chips,omitempty"`
	Item  *Item  `json:"item,omitempty"`
}
type Item struct {
	ItemID string `json:"itemId"`
	Days   int64  `json:"days"`
}

func ToWoF(wof *domain.WheelOfFortune) *WoF {
	return &WoF{
		Slots:      toSlots(wof.Slots),
		WonSlotIdx: wof.WonSlotIdx,
	}
}

func toSlots(s []domain.Slot) []Slot {
	slots := make([]Slot, 0, len(s))
	for _, slot := range s {
		slots = append(slots, toSlot(slot))
	}
	return slots
}

func toSlot(s domain.Slot) Slot {
	if s.Item != nil {
		return Slot{Item: &Item{
			ItemID: s.Item.ItemID,
			Days:   s.Item.Days,
		}}
	}
	return Slot{Chips: &s.Chips}
}

var emptyWoF = WoF{
	Slots:      []Slot{},
	WonSlotIdx: -1,
}

func PublishEmptyWheel(s *domain.Survival) error {
	if conf.IsE2E() {
		return nil
	}
	err := mqws.PublishNews(&mqws.Message{ToUserIds: []string{s.UserID}, Events: []*mqws.Event{{Type: "survivalWheel", Payload: gomq.M{"wheel": emptyWoF}.JSON()}}})
	if err != nil {
		log.Errorf("failed to publish empty survivalWheel event: %s", err)
		return err
	}
	return nil
}

func PublishWheel(s *domain.Survival) error {
	if conf.IsE2EVar {
		return nil
	}
	err := mqws.PublishNews(&mqws.Message{ToUserIds: []string{s.UserID}, Events: []*mqws.Event{{Type: "survivalWheel", Payload: gomq.M{"wheel": ToWoF(s.GetWheelOfFortune())}.JSON()}}})
	if err != nil {
		log.Errorf("failed to publish survivalWheel event: %s", err)
		return err
	}
	return nil
}
