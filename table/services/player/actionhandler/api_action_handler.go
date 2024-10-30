package actionhandler

import (
	"context"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/broadcast"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/services/player/playerbank"
	"go.mongodb.org/mongo-driver/bson"
)

type DbUpdates func(*domain.Table) []bson.E
type WsEvents func(*domain.Table) []*events.TableEvent
type Timeout func(*domain.Table)

type ActionHandler interface {
	DbUpdates(*domain.Table) []bson.E
	WsEvents(*domain.Table) []*events.TableEvent
	Timeout(*domain.Table)
}

func HandleNoCtx(table *domain.Table) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DefaultTimeout)
	defer cancel()
	return Handle(ctx, table)
}

func Handle(ctx context.Context, table *domain.Table) error {
	handler := GetActionHandler(table)

	err := db.SetTableGameFlow(ctx, table.ID, table.GameFlowVersion, handler.DbUpdates(table))
	if err != nil {
		return err
	}

	tableEvents := events.BuildActionsOfPlayers(table)
	tableEvents = append(tableEvents, handler.WsEvents(table)...)
	tableEvents = append(tableEvents, playerbank.HandleNullifiedPlayersLeft(table)...)
	broadcast.SendTableEvents(table.ID.Hex(), tableEvents)

	var userEvents []events.UserEvents
	for _, p := range table.PlayersWithUpgradedIntents() {
		userEvents = append(userEvents, events.UserEvents{UserID: p.UserId, Events: []*events.TableEvent{events.BuildIntent(p)}})
	}
	if len(userEvents) > 0 {
		broadcast.SendTableEventsToUsers(table.ID.Hex(), userEvents, nil, nil, nil)
	}

	handler.Timeout(table)

	return nil
}
