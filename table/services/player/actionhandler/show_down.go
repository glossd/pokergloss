package actionhandler

import (
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/events"
	"go.mongodb.org/mongo-driver/bson"
)

type ShowDown struct{}

func (s ShowDown) DbUpdates(table *domain.Table) []bson.E {
	return db.AllTableUpdatesGameFlow(table)
}

func (s ShowDown) WsEvents(table *domain.Table) []*events.TableEvent {
	var es []*events.TableEvent
	if sofp := table.BuildStackOverflowPlayer(); sofp != nil {
		es = append(es, events.BuildStackOverflowPlayer(sofp))
	}
	for _, p := range table.ShowedDownPlayers() {
		es = append(es, events.BuildShowDown(p, table.DecidingPosition))
	} /**/
	es = append(es, events.TimeToDecideBuilder(table, true))
	return es
}

func (s ShowDown) Timeout(table *domain.Table) {
	LaunchDecisionTimeout(table)
}
