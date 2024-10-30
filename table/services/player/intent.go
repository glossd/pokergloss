package player

import (
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/broadcast"
	"github.com/glossd/pokergloss/table/services/events"
	"go.mongodb.org/mongo-driver/bson"
)

func SetIntent(params IntentParams) error {
	table, err := db.FindTable(params.ctx, params.tableID)
	if err != nil {
		return err
	}

	err = table.SetIntent(params.position, params.iden, domain.NewIntent(params.Intent.Type, params.Intent.Chips))
	if err != nil {
		return err
	}

	player := table.GetPlayerUnsafe(params.position)

	err = db.SetTableFilter(params.ctx, table.ID,
		[]bson.E{{Key: "decidingposition", Value: bson.M{"$ne": player.Position}}},
		[]bson.E{{Key: db.PlayerDbPath(player.Position) + ".intent", Value: player.Intent}})
	if err != nil {
		return err
	}

	sendIntentToUser(table.ID.Hex(), player)

	return nil
}

func RemoveIntent(params *PositionParams) error {
	table, err := db.FindTable(params.ctx, params.tableID)
	if err != nil {
		return err
	}

	err = table.RemoveIntent(params.position, params.iden)
	if err != nil {
		return err
	}

	err = db.SetTableContext(params.ctx, table.ID, []bson.E{
		{Key: db.PlayerDbPath(params.position) + ".intent", Value: nil},
	})
	if err != nil {
		return err
	}

	sendIntentToUser(table.ID.Hex(), table.GetPlayerUnsafe(params.position))

	return nil
}

func sendIntentToUser(tableID string, player *domain.Player) {
	ue := []events.UserEvents{{UserID: player.UserId, Events: []*events.TableEvent{events.BuildIntent(player)}}}
	broadcast.SendTableEventsToUsers(tableID, ue, nil, nil, nil)
}
