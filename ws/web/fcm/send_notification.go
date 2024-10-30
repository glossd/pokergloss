package fcm

import (
	"context"
	"firebase.google.com/go/messaging"
	"fmt"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/ws/db"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SendNotificationTo(userID string, events []*mqws.Event) {
	if !conf.IsProd() {
		return
	}

	for _, event := range events {
		log.Debugf("Got event type=%s", event.Type)
		fmcNotification, webconfig := buildNotification(event)
		if fmcNotification != nil {
			log.Debugf("Have built notification fot type=%s", event.Type)
			ctx := context.Background()
			n, err := db.FindNotification(ctx, userID)
			if err != nil {
				if err != mongo.ErrNoDocuments {
					log.Errorf("Failed to send notification: %s", err)
				}
				continue
			}

			fcmMsg := &messaging.Message{
				Notification: fmcNotification,
				Token:        n.Web.Token,
				Webpush:      webconfig,
			}
			log.Debugf("Sending message to FCM, type=%s", event.Type)
			_, err = FcmClient.Send(ctx, fcmMsg)
			if err != nil {
				log.Errorf("FCM client failed to send a msg: %s", err)
				return
			}
		}
	}
}

func buildNotification(e *mqws.Event) (*messaging.Notification, *messaging.WebpushConfig) {
	switch e.Type {
	case "newsMessengerNewMessage":
		text := gjson.Get(e.Payload, "text").String()
		username := gjson.Get(e.Payload, "from.username").String()
		title := fmt.Sprintf("Message from %s", username)
		return &messaging.Notification{
				Title:    title,
				Body:     text,
				ImageURL: conf.Props.LogoURL,
			}, &messaging.WebpushConfig{
				Notification: &messaging.WebpushNotification{
					Title:              title,
					Body:               text,
					Icon:               conf.Props.LogoURL,
					RequireInteraction: false,
				},
				FcmOptions: &messaging.WebpushFcmOptions{
					Link: fmt.Sprintf("https://%s/messenger", conf.Props.Domain),
				},
			}
	case "multiGameStart":
		name := gjson.Get(e.Payload, "name").String()
		tableId := gjson.Get(e.Payload, "tableId").String()
		title := fmt.Sprintf("Tournament %s just started", name)
		text := "Click to join your table"
		return &messaging.Notification{
				Title:    title,
				Body:     text,
				ImageURL: conf.Props.LogoURL,
			}, &messaging.WebpushConfig{
				Notification: &messaging.WebpushNotification{
					Title:              title,
					Body:               text,
					Icon:               conf.Props.LogoURL,
					RequireInteraction: true,
				},
				FcmOptions: &messaging.WebpushFcmOptions{
					Link: fmt.Sprintf("https://%s/tables/%s", conf.Props.Domain, tableId),
				},
			}
	case "sitngoGameStart":
		name := gjson.Get(e.Payload, "name").String()
		tableId := gjson.Get(e.Payload, "tableId").String()
		title := fmt.Sprintf("SitNGo %s just started", name)
		text := "Click to join your table"
		return &messaging.Notification{
				Title:    title,
				Body:     text,
				ImageURL: conf.Props.LogoURL,
			}, &messaging.WebpushConfig{
				Notification: &messaging.WebpushNotification{
					Title:              title,
					Body:               text,
					Icon:               conf.Props.LogoURL,
					RequireInteraction: true,
				},
				FcmOptions: &messaging.WebpushFcmOptions{
					Link: fmt.Sprintf("https://%s/tables/%s", conf.Props.Domain, tableId),
				},
			}
	}
	return nil, nil
}
