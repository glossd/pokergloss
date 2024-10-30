package multi

import (
	"context"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/web/client/mq"
	"github.com/glossd/pokergloss/table/web/client/mqpub"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CountDownRebalance(ctx context.Context, event *mq.MultiPlayersMovedEvent) error {
	loid, err := primitive.ObjectIDFromHex(event.LobbyID)
	if err != nil {
		log.Errorf("Failed to parse lobbyID=%s : %s", event.LobbyID, err)
		return nil
	}

	// todo DG changed msgID to loid
	rc, err := db.FindRebalanceConfigAndCountDown(ctx, loid.String(), loid)
	if err != nil {
		log.Errorf("Failed to find rebalance config: %s", err)
		return err
	}
	defer func() {
		db.RecursiveUnlockRebalanceConfig(loid)
	}()

	if rc.CountDown > 0 {
		return nil
	}

	result, err := Rebalance(loid)
	if err != nil {
		return err
	}
	log.Infof("Rebalance result: %+v", result)
	rc = &db.RebalancerConfig{
		LobbyID:   loid,
		CountDown: result.CountTablesWithMovingPlayers,
	}
	switch result.Status {
	case StopRebalance:
		return nil
	case RemoveRightAway:
		mqpub.PublishMultiRebalanceAt(event.LobbyID, timeutil.NowAdd(conf.Props.RebalancerPeriod))
		return nil
	case MoveAllPlayers:
		_ = db.UpdateRebalanceConfig(ctx, rc)
		// will send rebalance event when all players are moved
	case Disproportion:
		if result.CountTablesWithMovingPlayers == 0 {
			mqpub.PublishMultiRebalanceAt(event.LobbyID, timeutil.NowAdd(conf.Props.RebalancerPeriod))
		} else {
			_ = db.UpdateRebalanceConfig(ctx, rc)
			// will send rebalance event on each table where tables are moved.
		}
	}
	return nil
}
