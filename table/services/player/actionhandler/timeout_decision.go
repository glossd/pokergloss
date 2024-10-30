package actionhandler

import (
	"context"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	"github.com/glossd/pokergloss/table/web/client/mqpub"
	log "github.com/sirupsen/logrus"
	"time"
)

var ErrTimeout = services.ErrFormat("time out")

// timeoutAt in milliseconds
func LaunchDecisionTimeout(table *domain.Table) {
	key := timeout.Key{TableID: table.ID, Position: table.DecidingPosition, Version: table.GameFlowVersion + 1}
	if table.DecisionTimeout == 1 {
		// for tests
		DoDecisionTimeoutNoCtx(key)
		return
	}
	if table.DecisionTimeout < 0 {
		// make it yourself. For tests
		return
	}
	_, err := table.DecidingPlayer()
	if err != nil {
		log.Panicf("LaunchDecisionTimeout, no deciding player, round=%s, table=%+v", table.RoundType(), table)
		return
	}

	mqpub.PublishTimeoutEvent(&timeout.Event{
		Type: timeout.Decision,
		At:   table.DecisionTimeoutAt,
		Key:  key,
	})
}

// For tests
func DoDecisionTimeoutNoCtx(key timeout.Key) (tryAgain bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return DoDecisionTimeout(ctx, key)
}

func DoDecisionTimeout(ctx context.Context, key timeout.Key) (tryAgain bool) {
	log.Debugf("Player decision timeout, key=%s", key)
	table, err := db.FindTableGameFlow(ctx, key.TableID, key.Version)
	if err != nil {
		if err == db.ErrVersionNotMatch {
			log.Tracef("Decision timeout race condition, key=%s: %s", key, err)
			return false
		}
		log.Errorf("Couldn't send timeToDecide timeout, find failed, key=%s : %s", key, err)
		return true
	}
	player, err := table.GetPlayer(key.Position)
	if err != nil {
		log.Errorf("Couldn't send timeToDecide timeout, getting stalePlayer by position, key=%s : %s", key, err)
		return false
	}

	err = table.MakeActionOnTimeout(key.Position)
	if err != nil {
		log.Errorf("Couldn't send timeToDecide timeout, making action on timeout, tableID=%s, player=%v : %s", table.ID, player, err)
		return false
	}

	err = Handle(ctx, table)
	if err != nil {
		if err == db.ErrVersionNotMatch {
			log.Tracef("Decision timeout race condition on handling, key=%s: %s", key, err)
			return false
		}
		log.Errorf("Couldn't send timeToDecide timeout, handle failed, tableID=%s, player=%v : %s", table.ID, player, err)
		return true
	}

	return
}
