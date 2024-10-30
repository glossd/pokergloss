package actionhandler

import (
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/web/client/mqpub"
	"go.mongodb.org/mongo-driver/bson"
)

type GameEnd struct{}

func (g GameEnd) DbUpdates(table *domain.Table) []bson.E {
	return db.AllTableUpdatesGameFlow(table)
}

func (g GameEnd) WsEvents(table *domain.Table) []*events.TableEvent {
	var es []*events.TableEvent
	if sofp := table.BuildStackOverflowPlayer(); sofp != nil {
		es = append(es, events.BuildStackOverflowPlayer(sofp))
	}
	for _, p := range table.ShowedDownPlayers() {
		es = append(es, events.BuildShowDown(p, table.DecidingPosition))
	}
	es = append(es, events.BuildWinners(table))
	return es
}

func (g GameEnd) Timeout(table *domain.Table) {
	mqpub.PublishGameEnd(table)
	LaunchDelayedGame(table)
}
