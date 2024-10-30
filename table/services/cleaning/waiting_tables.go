package cleaning

import (
	"context"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/broadcast"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/services/player/playerbank"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func CleanWaitingTables() int {
	var count int
	err := db.ForEachTable(bson.D{{Key: "type", Value: domain.CashType}, {Key: "status", Value: domain.WaitingTable}, {Key: "ispersistent", Value: false}}, func(ctx context.Context, t *domain.Table) {
		if timeutil.NowMinus(t.WaitingAt) > conf.Props.Cleaning.WaitingTablesTimeout {
			var acc []*events.TableEvent
			for _, p := range t.AllPlayers() {
				playerbank.SendPlayerChipsToBank(p, t)
				acc = append(acc, events.BuildPlayerLeft(p))
			}
			broadcast.SendTableEvents(t.ID.Hex(), acc)
			_ = db.DeleteTableNoCtx(t.ID)
			count++
		}
	})
	if err != nil {
		log.Errorf("Failed to clean waiting tables: %s", err)
	}
	if count > 0 {
		log.Infof("Deleted waiting tables, count=%d", count)
	}
	return count
}
