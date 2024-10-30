package rpcstore

import (
	"github.com/glossd/pokergloss/gogrpc/grpctableevents"
	"github.com/glossd/pokergloss/gomq/mqws"
	"sync"
)

var theMap = &sync.Map{}

type Sender interface {
	Send(*grpctableevents.Events) error
}

func AddSender(tableID string, s Sender) {
	theMap.Store(tableID, s)
}

func RemoveSender(tableID string) {
	theMap.Delete(tableID)
}

func SendEvents(tableID string, msg *mqws.TableMessage) error {
	v, ok := theMap.Load(tableID)
	if ok {
		sender := v.(Sender)
		err := sender.Send(&grpctableevents.Events{Events: ExtractEvents(msg)})
		if err != nil {
			return err
		}
	}
	return nil
}

func ExtractEvents(msg *mqws.TableMessage) []*grpctableevents.Event {
	ue := msg.UserEvents
	if ue != nil {
		events := make([]*grpctableevents.Event, 0, len(ue.BeforeEvents)+len(ue.SecretEvents)+len(ue.BeforeEvents))
		events = append(events, toGrpcEvents(ue.BeforeEvents)...)
		events = append(events, toGrpcEvents(ue.SecretEvents)...)
		events = append(events, toGrpcEvents(ue.AfterEvents)...)
		return events
	} else {
		return toGrpcEvents(msg.Events)
	}
}

func toGrpcEvents(a []*mqws.Event) []*grpctableevents.Event {
	result := make([]*grpctableevents.Event, 0, len(a))
	for _, e := range a {
		result = append(result, &grpctableevents.Event{Type: e.Type, Payload: e.Payload})
	}
	return result
}
