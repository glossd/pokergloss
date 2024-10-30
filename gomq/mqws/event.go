package mqws

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type event struct {
	Type string `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

func (x *Message) EventsToJson() []byte {
	return EventsToJson(x.Events)
}

func EventsToJson(es []*Event) []byte {
	events := make([]*event, 0, len(es))
	for _, e := range es {
		var payload map[string]interface{}
		err := json.Unmarshal([]byte(e.Payload), &payload)
		if err != nil {
			log.Errorf("mqws.EventsToJson couldn't parse event payload to json, payload=%s: %s", e.Payload, err)
			continue
		}
		events = append(events, &event{Type: e.Type, Payload: payload})
	}
	marshal, err := json.Marshal(events)
	if err != nil {
		log.Errorf("mqws.EventsToJson couldn't marshal evennts to json ")
		return nil
	}
	return marshal
}
