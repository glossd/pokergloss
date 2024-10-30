package mq

import (
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/ws/storage"
	"github.com/glossd/pokergloss/ws/web/fcm"
	log "github.com/sirupsen/logrus"
)

func SendMessageToWs(msg *mqws.Message) {
	log.Debugf("Got message from mq")
	if len(msg.ToUserIds) == 0 {
		if msg.EntityId != "" {
			send(msg, msg.EntityId)
		}
	} else {
		for _, userId := range msg.ToUserIds {
			send(msg, userId)
		}
	}
}

func send(msg *mqws.Message, userId string) {
	hub, ok := storage.GetUserHub(userId)
	log.Debugf("User message, user found %t", ok)
	if ok {
		hub.Broadcast <- msg.EventsToJson()
		if hub.IsZeroUserConn() {
			fcm.SendNotificationTo(userId, msg.Events)
		}
	} else {
		fcm.SendNotificationTo(userId, msg.Events)
	}
}
