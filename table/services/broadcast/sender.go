package broadcast

import (
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/web/client/mqpub"
)

func SendManyTableEvent(tableIDs []string, tableEvents []*events.TableEvent) {
	mqpub.SendManyTableMessage(tableIDs, tableEvents)
}

func SendTableEvents(tableID string, tableEvents []*events.TableEvent) {
	mqpub.SendTableMessage(tableID, tableEvents)
}

func SendTableEvent(tableID string, event *events.TableEvent) {
	SendTableEvents(tableID, []*events.TableEvent{event})
}

func SendTableEventsToUsers(tableID string, ue []events.UserEvents, notFoundUserEvents, beforeEvents, afterEvents []*events.TableEvent) {
	mqpub.SendTableMessageToUsers(tableID, ue, notFoundUserEvents, nil, beforeEvents, afterEvents)
}
