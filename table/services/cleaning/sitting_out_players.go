package cleaning

import (
	"context"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/broadcast"
	"github.com/glossd/pokergloss/table/services/player/playerbank"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func CleanSittingOutPlayers() {
	err := db.ForEachTable(bson.D{}, removeSittingOutPlayer)
	if err != nil {
		log.Errorf("Failed to clean of sitting out players : %s", err)
	}
}

func removeSittingOutPlayer(ctx context.Context, table *domain.Table) {
	playersToNullifyCandidates := table.PlayersFilter(func(p *domain.Player) bool {
		var timeout time.Duration
		if table.IsCashType() {
			timeout = conf.Props.Cleaning.CashSittingOutTimeout
		} else {
			timeout = conf.Props.Cleaning.TournamentSittingOutTimeout
		}

		return p.IsSittingOut() && timeutil.NowMinus(p.SatOutAt) > timeout
	})
	if len(playersToNullifyCandidates) == 0 {
		return
	}

	playersToNullify := make([]*domain.Player, 0, len(playersToNullifyCandidates))
	for _, p := range playersToNullifyCandidates {
		err := table.NullifySittingOutPlayer(p.Position)
		if err != nil {
			log.Errorf("CleanSittingOutPlayers: %s", err)
			continue
		}
		playersToNullify = append(playersToNullify, p)
	}

	dbUpdates := make([]bson.E, 0, len(playersToNullifyCandidates))
	for _, player := range playersToNullifyCandidates {
		dbUpdates = append(dbUpdates, db.PlayerNullify(player.Position))
	}
	dbUpdates = append(dbUpdates, db.PlayersCount(table))

	err := db.SetTable(table.ID, dbUpdates)
	if err != nil {
		log.Errorf("Couldn't remove sitting out players, tableID=%s : %s", table.ID.Hex(), err)
		return
	}

	wsEvents := playerbank.HandleNullifiedPlayersLeft(table)
	broadcast.SendTableEvents(table.ID.Hex(), wsEvents)
}
