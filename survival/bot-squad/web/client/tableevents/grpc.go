package tableevents

import (
	"github.com/glossd/pokergloss/gogrpc/grpctableevents"
	"github.com/glossd/pokergloss/survival/bot-squad/conf"
	"github.com/glossd/pokergloss/survival/bot-squad/domain"
	"github.com/glossd/pokergloss/survival/bot-squad/service"
	"github.com/glossd/pokergloss/table-stream/web/grpc"
	log "github.com/sirupsen/logrus"
)

func StreamEventsGRPC(c conf.Config) {
	err := grpc.StreamTableEventsLocal(c.TableID, func(msg *grpctableevents.Events) {
		events := make([]*domain.Event, 0, len(msg.Events))
		for _, event := range msg.Events {
			events = append(events, domain.NewEvent(event.Type, event.Payload))
		}
		service.HandleEvents(c, events)
	})
	if err != nil {
		log.Errorf("Failed to stream events: %s", err)
	}
}
