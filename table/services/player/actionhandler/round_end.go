package actionhandler

import (
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/events"
	"go.mongodb.org/mongo-driver/bson"
)

type RoundEnd struct{}

func (r RoundEnd) DbUpdates(table *domain.Table) []bson.E {
	return db.AllTableUpdatesGameFlow(table)
}

func (r RoundEnd) WsEvents(table *domain.Table) []*events.TableEvent {
	var tableEvents []*events.TableEvent
	tableEvents = append(tableEvents, events.BuildNewBettingRound(table))
	tableEvents = append(tableEvents, events.BuildTimeToDecide(table))
	return tableEvents
}

func (r RoundEnd) Timeout(table *domain.Table) {
	LaunchDecisionTimeout(table)
}
