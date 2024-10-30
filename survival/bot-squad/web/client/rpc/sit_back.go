package rpc

import (
	"context"
	"github.com/glossd/pokergloss/gogrpc/grpctable"
	"github.com/glossd/pokergloss/survival/bot-squad/conf"
	grpc "github.com/glossd/pokergloss/table/web/grpcserver"
	log "github.com/sirupsen/logrus"
)

func SitBack(config conf.Config, position int) {
	_, err := grpc.SitBack(context.Background(), &grpctable.SitBackRequest{
		TableId:  config.TableID,
		Position: int64(position),
	})
	if err != nil {
		log.Errorf("Sit back failed %v", err)
	}
}
