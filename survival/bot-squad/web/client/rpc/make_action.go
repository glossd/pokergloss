package rpc

import (
	"context"
	"github.com/glossd/pokergloss/gogrpc/grpctable"
	"github.com/glossd/pokergloss/survival/bot-squad/conf"
	"github.com/glossd/pokergloss/survival/bot-squad/domain"
	grpc "github.com/glossd/pokergloss/table/web/grpcserver"
	log "github.com/sirupsen/logrus"
)

func MakeAction(config conf.Config, position int, a domain.Action, t *domain.Table) {

	_, err := grpc.MakeAction(context.Background(), &grpctable.MakeActionRequest{
		TableId:    config.TableID,
		Position:   int64(position),
		ActionType: string(a.Type),
		Chips:      a.Chips,
	})
	if err != nil {
		log.Errorf("Make action %v, maxRoundBet=%d, totalRoundBet=%d failed %v", a, t.MaxRoundBet, t.DecidingPlayer().TotalRoundBet, err)
	}
}
