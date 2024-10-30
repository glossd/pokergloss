package player

import (
	"errors"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services"
	"github.com/glossd/pokergloss/table/services/broadcast"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrSitBackRaceCondition = services.ErrFormat("state of your table has changed too fast, please sit back again")

func SitBack(params *PositionParams) (*domain.Player, error) {
	table, err := db.FindTable(params.ctx, params.tableID)
	if err != nil {
		return nil, err
	}

	isNewGame, err := table.SitBack(params.position, params.iden)
	if err != nil {
		if err == domain.ErrNoSitBackInSittingOut {
			return table.GetPlayerUnsafe(params.position), err
		}
		return nil, err
	}

	var dbUpdates []primitive.E
	seatPlayerBack := table.GetSeatUnsafe(params.position)
	if isNewGame { // mongo can't update the same elements
		// there are only two playing players
		dbUpdates = db.AllTableUpdatesGameFlow(table)
	} else {
		dbUpdates = db.SeatUpdate(seatPlayerBack)
	}

	var dbErr error
	if isNewGame {
		dbErr = db.SetTableGameFlow(params.ctx, table.ID, table.GameFlowVersion, dbUpdates)
	} else {
		dbErr = db.SetTableContext(params.ctx, table.ID, dbUpdates)
	}
	if dbErr != nil {
		if errors.Is(dbErr, db.ErrVersionNotMatch) {
			return nil, ErrSitBackRaceCondition
		}
		log.Errorf("Couldn't sit back %s: %s", params.iden, err)
		return nil, err
	}

	broadcast.SendTableEvent(params.tableID.Hex(), events.BuildPlayerSitBack(seatPlayerBack.GetPlayer()))

	if isNewGame {
		actionhandler.HandleGameStart(table)
	}

	return nil, nil
}
