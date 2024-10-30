package mq

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq/mqprofile"
	"github.com/glossd/pokergloss/profile/domain"
	log "github.com/sirupsen/logrus"
)

func PublishProfileUpdate(doc *domain.Profile) {
	PublishProfileUpdateFields(doc.UserID, doc.Username, doc.Picture)
}

func PublishCreated(p *domain.Profile) {
	if conf.IsProd() {
		err := mqprofile.PublishCreated(&mqprofile.Profile{
			UserId:   p.UserID,
			Username: p.Username,
			Picture:  p.Picture,
		})
		if err != nil {
			log.Errorf("Failed to publish created profile: %s", err)
		}
	}
}

func PublishProfileUpdateFields(userID, username, picture string) {
	if conf.IsProd() {
		err := mqprofile.Publish(&mqprofile.Profile{
			UserId:   userID,
			Username: username,
			Picture:  picture,
		})
		if err != nil {
			log.Errorf("Failed to publish profile: %s", err)
		}
	}
}
