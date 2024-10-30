package ws

import (
	"github.com/glossd/pokergloss/achievement/domain"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/gomq/mqws"
	log "github.com/sirupsen/logrus"
)

type AchieveEvent struct {
	Type    string      `json:"type" enums:"expUpdate,hand,newAchievement"`
	Payload interface{} `json:"payload"`
}

// @ID events from websocket
// @Success 200 {array} AchieveEvent
// @Router /ws [get]
func PublishExpEvent(exp *domain.ExP) {
	if !conf.IsProd() {
		return
	}
	err := mqws.Publish(&mqws.Message{
		EntityType: mqws.Message_USER,
		EntityId:   exp.UserID,
		Events:     []*mqws.Event{{Type: "expUpdate", Payload: gomq.M{"level": exp.Level, "points": exp.Points, "nextLevelPoints": exp.NextLevelPoints(), "startLevelPoints": exp.StartLevelPoints(), "prize": exp.GetNewLevelPrize()}.JSON()}},
	})
	if err != nil {
		log.Errorf("Failed to publish exp event: %s", err)
	}
}

func PublishHandEvent(as *domain.AchievementStore) {
	if !conf.IsProd() {
		return
	}
	hc := as.HandsCounter
	if hc.GetPrize().Chips == 0 {
		return
	}
	err := mqws.Publish(&mqws.Message{
		EntityType: mqws.Message_USER,
		EntityId:   as.UserID,
		Events: []*mqws.Event{{Type: "hand", Payload: gomq.M{
			"prize": gomq.M{
				"hand":  hc.GetPrize().Name,
				"name":  hc.GetPrize().Name,
				"type":  domain.HandType,
				"chips": hc.GetPrize().Chips,
				"count": hc.GetPrizeHandCount(),
			}}.JSON()}},
	})
	if err != nil {
		log.Errorf("Failed to publish hand event: %s", err)
	}
}

func PublishNewCounterAchievement(userID string, c domain.Counter) {
	if !conf.IsProd() {
		return
	}
	if c.GetPrize().Chips == 0 {
		return
	}
	err := mqws.Publish(&mqws.Message{
		EntityType: mqws.Message_USER,
		EntityId:   userID,
		Events: []*mqws.Event{{Type: "newAchievement", Payload: gomq.M{
			"prize": gomq.M{
				"name":  c.GetName(),
				"chips": c.GetPrize().Chips,
				"type":  c.GetType(),
				"count": c.GetCount(),
			}}.JSON()}},
	})
	if err != nil {
		log.Errorf("Failed to publish newAchievement event: %s", err)
	}
}
