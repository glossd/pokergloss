package mqpub

import (
	conf "github.com/glossd/pokergloss/goconf"
	mqtable2 "github.com/glossd/pokergloss/gomq/mqtable"
	"github.com/glossd/pokergloss/table/domain"
	log "github.com/sirupsen/logrus"
)

func PublishSurvivalEnd(t *domain.Table) {
	if conf.IsLocalOnly() {
		return
	}
	if !t.IsSurvival {
		return
	}

	if conf.IsProd() {
		var userID string
		if t.IsSurvivalUserLeft() {
			for _, p := range t.NullifiedLeavingPlayers() {
				if p.Position == 0 {
					userID = p.UserId
				}
			}
		} else {
			p, err := t.GetPlayer(0)
			if err == nil {
				userID = p.UserId
			}
		}
		err := mqtable2.PublishSurvivalEnd(&mqtable2.SurvivalEnd{
			TableId:    t.ID.Hex(),
			UserId:     userID,
			IsUserLost: t.IsSurvivalUserLeft()})
		if err != nil {
			log.Errorf("PublishSurvivalEnd failed: %s", err)
		}
	}
}
