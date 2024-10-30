package mqpub

import (
	"encoding/json"
	conf "github.com/glossd/pokergloss/goconf"
	mqws2 "github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/web/client/mq"
	log "github.com/sirupsen/logrus"
)

func SendTableMessage(tableID string, events []*events.TableEvent) {
	SendManyTableMessage([]string{tableID}, events)
}

func SendManyTableMessage(tableIDs []string, events []*events.TableEvent) {
	if len(events) == 0 {
		return
	}
	msg := &mqws2.TableMessage{
		ToEntityIds: tableIDs,
		Events:      ToEvents(events),
	}
	publishTableMsg(msg)
}

func SendTableMessageToUsers(tableID string, ues []events.UserEvents, notFoundUserEvents, secretEvents, beforeEvents, afterEvents []*events.TableEvent) {
	userEvents := make(map[string]*mqws2.Events)
	for _, ue := range ues {
		userEvents[ue.UserID] = &mqws2.Events{Events: ToEvents(ue.Events)}
	}

	msg := &mqws2.TableMessage{
		ToEntityIds: []string{tableID},
		UserEvents: &mqws2.TableUserEvents{
			UserEvents:          userEvents,
			NotFoundUsersEvents: ToEvents(notFoundUserEvents),
			SecretEvents:        ToEvents(secretEvents),
			BeforeEvents:        ToEvents(beforeEvents),
			AfterEvents:         ToEvents(afterEvents),
		},
	}
	publishTableMsg(msg)
}

func publishTableMsg(msg *mqws2.TableMessage) {
	if conf.IsProd() || conf.IsLocalOnly() {
		err := mqws2.PublishTableMsg(msg)
		if err != nil {
			log.Errorf("Failed to publish ws message: %v", err)
		} else {
			log.Tracef("Published mqws.Message msg: %v", msg)
		}
	} else if conf.IsE2E() {
		mq.TestMQ <- msg
	}
}

func ToEvents(events []*events.TableEvent) []*mqws2.Event {
	var result []*mqws2.Event
	for _, event := range events {
		if event != nil {
			eventToAdd := ToEvent(event)
			if eventToAdd != nil {
				result = append(result, eventToAdd)
			}
		}
	}
	return result
}

func ToEvent(e *events.TableEvent) *mqws2.Event {
	bytes, err := json.Marshal(e.Payload)
	if err != nil {
		log.Errorf("Failed to marshal event=%+v: %v", e, err)
		return nil
	}
	return &mqws2.Event{
		Type:    string(e.Type),
		Payload: string(bytes),
	}
}
