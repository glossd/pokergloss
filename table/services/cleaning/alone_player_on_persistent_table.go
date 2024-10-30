package cleaning

import (
	"context"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/player"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func CleanAlonePlayerOnPersistentTable() int {
	count := 0
	err := db.ForEachTable(bson.D{{Key: "type", Value: domain.CashType}, {Key: "status", Value: domain.WaitingTable}, {Key: "ispersistent", Value: true}, {Key: "playerscount", Value: 1}}, func(ctx context.Context, t *domain.Table) {
		if timeutil.NowMinus(t.WaitingAt) > conf.Props.Cleaning.AlonePlayerOnPersistentTable {
			if p, ok := t.IsOnlyOneOnTable(); ok {
				if p.Username == "pokerblow" {
					return
				}
				params, err := player.NewPositionParams(ctx, t.ID.Hex(), p.Position, p.Identity)
				if err != nil {
					log.Panicf("Alone player, positions params error %s", err)
				}
				err = player.Stand(params)
				if err != nil {
					log.Errorf("Failed to make player stand: %s", err)
					return
				}
				count++
			}
		}
	})
	if err != nil {
		log.Errorf("Failed to clean waiting tables: %s", err)
	}
	if count > 0 {
		log.Infof("Released alone players from persistent tables, count=%d", count)
	}
	return count
}
