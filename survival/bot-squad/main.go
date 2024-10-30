package botsquad

import (
	"github.com/glossd/pokergloss/survival/bot-squad/conf"
	"github.com/glossd/pokergloss/survival/bot-squad/service"
	"github.com/glossd/pokergloss/survival/bot-squad/web/client/tableevents"
)

func Create(c conf.Config) error {
	if c.Protocol != "ws" {
		// start streaming before creation of the table
		go tableevents.StreamEventsGRPC(c)
	}
	err := service.Init(c)
	if err != nil {
		return err
	}

	if c.Protocol == "ws" {
		go tableevents.StreamEventsWS(c)
	}
	return nil
}

func Delete(userID string) error {
	service.DeleteUserData(userID)
	return nil
}
