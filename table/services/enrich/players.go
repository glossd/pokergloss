package enrich

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/broadcast"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/services/playerautoconfig"
	"github.com/glossd/pokergloss/table/web/client/market"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func Players(staleTable *domain.Table, players []*domain.Player) {
	if !conf.Props.Enrich.PlayersEnabled {
		return
	}
	if len(players) == 0 {
		return
	}

	for _, p := range players {
		item, err := market.GetSelectedItemID(p.UserId)
		if err != nil {
			continue
		}
		p.SetMarketItem(item.ID, item.CoinsDayPrice)
	}

	for _, player := range players {
		pac, err := playerautoconfig.FindAutoConfigNoCtx(player.Identity)
		if err != nil {
			log.Errorf("Failed to enrich auto config: %s", err)
			continue
		}
		if pac != nil {
			staleTable.SetPlayerAutoConfig(player, *pac)
		}
	}

	var dbUpdates []bson.E
	for _, p := range players {
		dbUpdates = append(dbUpdates, db.PlayerEnrichment(p)...)
	}
	err := db.SetTable(staleTable.ID, dbUpdates)
	if err != nil {
		log.Errorf("Failed to enrich players: %s", err)
		return
	}

	broadcast.SendTableEvent(staleTable.ID.Hex(), events.BuildEnrichment(players))

	var userEvents []events.UserEvents
	for _, player := range players {
		userEvents = append(userEvents, events.UserEvents{
			UserID: player.UserId,
			Events: events.BuildSetPlayerConfig(player.GetAutoConfig()),
		})
	}
	broadcast.SendTableEventsToUsers(staleTable.ID.Hex(), userEvents, nil, nil, nil)
}
