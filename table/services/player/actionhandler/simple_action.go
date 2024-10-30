package actionhandler

import (
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/events"
	"go.mongodb.org/mongo-driver/bson"
)

type Action struct{}

func (a Action) DbUpdates(table *domain.Table) []bson.E {
	return db.AllTableUpdatesGameFlow(table)
}

func (a Action) WsEvents(table *domain.Table) []*events.TableEvent {
	return []*events.TableEvent{events.BuildTimeToDecide(table)}
}

func (a Action) Timeout(table *domain.Table) {
	LaunchDecisionTimeout(table)
}
