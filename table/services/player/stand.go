package player

import (
	"errors"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/broadcast"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	"github.com/glossd/pokergloss/table/services/player/playerbank"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var ErrStandRaceCondition = errors.New("state of your table has changed too fast, please stand again")

func Stand(params *PositionParams) error {
	ctx := params.ctx
	tableID := params.tableID
	iden := params.iden
	position := params.position

	table, err := db.FindTable(ctx, tableID)
	if err != nil {
		log.Errorf("User couldn't leave table, finding table error: %s", err)
		return err
	}

	player, err := table.Stand(position, iden)
	if err != nil {
		return err
	}

	if player.IsStandFolded() {
		err := actionhandler.HandleNoCtx(table)
		if err != nil {
			log.Errorf("Couldn't make auto-fold on stand, tableID=%s, position=%d : %s", tableID.Hex(), position, err)
		}
		return err
	}

	isSimpleStand := player.Status == domain.PlayerReservedSeat ||
		player.Status == domain.PlayerSittingOut ||
		(!table.IsGameEnd() && player.LastGameAction == domain.Fold) ||
		(!table.IsGameEnd() && table.IsEnoughPlayersForGame())

	var dbUpdates []bson.E
	var wsEvents []*events.TableEvent
	if table.IsSeatFree(position) {
		if isSimpleStand {
			dbUpdates = append(dbUpdates, db.PlayerNullify(player.Position))
		} else {
			dbUpdates = db.AllTableUpdatesGameFlow(table)
		}
		wsEvents = append(wsEvents, events.BuildPlayerLeft(player))
	} else {
		dbUpdates = append(dbUpdates, db.PlayerLeaving(position))
	}

	var dbErr error
	if table.IsSeatFree(position) && !isSimpleStand {
		// stand in the middle of game end time case
		dbErr = db.SetTableGameFlow(params.ctx, table.ID, table.GameFlowVersion, dbUpdates)
	} else {
		dbErr = db.SetTableContext(params.ctx, tableID, dbUpdates)
	}

	if dbErr != nil {
		if errors.Is(dbErr, db.ErrVersionNotMatch) {
			return ErrStandRaceCondition
		} else {
			log.Errorf("User couldn't stand from table, userID=%s, tableID=%s: %s", iden.UserId, tableID.Hex(), dbErr)
			return dbErr
		}
	}

	if table.IsSeatFree(position) {
		playerbank.SendPlayerChipsToBank(player, table)
	}

	broadcast.SendTableEvents(tableID.Hex(), wsEvents)
	return nil
}
