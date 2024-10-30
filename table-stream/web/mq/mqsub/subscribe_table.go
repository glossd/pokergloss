package mqsub

import (
	"context"
	"github.com/glossd/memmq"
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/table-stream/web/ws"
	log "github.com/sirupsen/logrus"
)

const SubID = "table-events-service"

func SubscribeTable() {
	err := mqws.SubscribeTableMsg(SubID, func(ctx context.Context, msg *mqws.TableMessage) error {
		log.Tracef("Message to %s", msg.ToEntityIds)
		for _, id := range msg.ToEntityIds { // there's gonna be ONLY ONE entity id
			if msg.UserEvents == nil {
				ws.Broadcast(id, mqws.EventsToJson(msg.Events))
			} else {
				ws.Direct(id, msg.UserEvents)
			}
			memmq.Publish("pg.table-stream.grpc."+id, msg)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("mqsub.SubscribeTable failed: %v", err)
	}
	return
}
