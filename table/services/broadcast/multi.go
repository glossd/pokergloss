package broadcast

import (
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/events"
	log "github.com/sirupsen/logrus"
)

func SendEventsAboutMovingPlayers(tableFrom *domain.Table, tablesTo []*domain.Table) {
	tellPlayersToMove(tableFrom, tablesTo)
	tellPlayersFromTablesToAboutArrivedPlayers(tablesTo)
}

func tellPlayersFromTablesToAboutArrivedPlayers(tablesTo []*domain.Table) {
	for _, t := range tablesTo {
		toEvents := make([]*events.TableEvent, 0, len(t.GetPutPlayers()))
		for _, player := range t.GetPutPlayers() {
			toEvents = append(toEvents, events.BuildBankroll(player))
		}
		SendTableEvents(t.ID.Hex(), toEvents)
	}
}

// Also tells players of tableFrom who left
func tellPlayersToMove(tableFrom *domain.Table, tablesTo []*domain.Table) {
	var ue []events.UserEvents
	var afterEvents []*events.TableEvent
	for _, t := range tablesTo {
		for _, player := range t.MultiAttrs.GetPutPlayers() {
			ue = append(ue, events.UserEvents{
				UserID: player.UserId,
				Events: []*events.TableEvent{events.BuildMultiPlayerMove(t.ID)},
			})
			afterEvents = append(afterEvents, events.BuildPlayerMoved(player))
		}
	}
	SendTableEventsToUsers(tableFrom.ID.Hex(), ue, nil, nil, afterEvents)
}

// needs to be sent after each game of each table, to update stacks of players
func SendMultiTablePlayersUpdates(table *domain.Table) {
	lobby, err := db.FindLobbyMultiNoCtx(table.LobbyID)
	if err != nil {
		log.Errorf("Failed multi to send player updates: %s", err)
		return
	}

	var es []*events.TableEvent
	if table.IsWaiting() && table.IsOnlyOneOnTableBool() {
		es = append(es, events.BuildMultiPlayersEmpty(table.ID))
	} else {
		es = append(es, events.BuildMultiPlayersUpdate(table))
	}
	ppt := table.MovedPlayersPerTable()
	for tableID, movedPlayer := range ppt {
		es = append(es, events.BuildMultiPlusPlayersUpdate(tableID, movedPlayer))
	}

	SendManyTableEvent(lobby.GetTableIDsAsStr(), es)
}
