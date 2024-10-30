package grpc

import (
	"github.com/glossd/memmq"
	"github.com/glossd/pokergloss/gogrpc/grpctableevents"
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/table-stream/web/grpc/rpcstore"
	log "github.com/sirupsen/logrus"
)

func StreamTableEvents(r *grpctableevents.StreamTableEventsRequest, stream grpctableevents.TableEventsService_StreamTableEventsServer) error {
	rpcstore.AddSender(r.TableId, stream)
	select {
	case <-stream.Context().Done():
		rpcstore.RemoveSender(r.TableId)
	}
	return nil
}

// Blocking!
func StreamTableEventsLocal(tableID string, process func(event *grpctableevents.Events)) error {
	return memmq.Subscribe("pg.table-stream.grpc."+tableID, func(msg interface{}) bool {
		v, ok := msg.(*mqws.TableMessage)
		if !ok {
			log.Errorf("Stram table events, expected *mqws.TableMessage, got %T", v)
			return true
		}
		process(&grpctableevents.Events{Events: rpcstore.ExtractEvents(v)})
		return true
	})
}
