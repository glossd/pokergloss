package domain

import (
	"github.com/glossd/pokergloss/gomq/mqws"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

type Event struct {
	Type      string
	TableIDs  []string
	CreatedAt bsonx.Val
	Payload   interface{}
}

func ToEvents(msg *mqws.TableMessage) []*Event {
	if len(msg.Events) > 0 {
		var result = make([]*Event, 0, len(msg.Events))
		for _, event := range msg.Events {
			e, err := ToEvent(msg.ToEntityIds, event)
			if err != nil {
				continue
			}
			result = append(result, e)
		}
		return result
	}
	if msg.UserEvents != nil {
		var result []*Event
		for _, event := range msg.UserEvents.BeforeEvents {
			e, err := ToEvent(msg.ToEntityIds, event)
			if err != nil {
				continue
			}
			result = append(result, e)
		}
		for _, event := range msg.UserEvents.SecretEvents {
			e, err := ToEvent(msg.ToEntityIds, event)
			if err != nil {
				continue
			}
			result = append(result, e)
		}
		if len(msg.UserEvents.SecretEvents) == 0 {
			for _, events := range msg.UserEvents.UserEvents {
				for _, event := range events.Events {
					e, err := ToEvent(msg.ToEntityIds, event)
					if err != nil {
						continue
					}
					result = append(result, e)
				}
			}
		}
		for _, event := range msg.UserEvents.AfterEvents {
			e, err := ToEvent(msg.ToEntityIds, event)
			if err != nil {
				continue
			}
			result = append(result, e)
		}
		return result
	}
	return nil
}

func ToEvent(tableIDs []string, e *mqws.Event) (*Event, error) {
	var payload interface{}
	err := bson.UnmarshalExtJSON([]byte(e.Payload), true, &payload)
	if err != nil {
		log.Errorf("Unmarshalling JSON error: %s", err)
		return nil, err
	}
	return &Event{
		Type:      e.Type,
		TableIDs:  tableIDs,
		CreatedAt: bsonx.DateTime(time.Now().UnixNano() / 1e6),
		Payload:   payload,
	}, nil
}
