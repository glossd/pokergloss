package player

import (
	"errors"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/services"
	"github.com/glossd/pokergloss/table/services/broadcast"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	log "github.com/sirupsen/logrus"
)

func ReserveTableSeat(params *PositionParams) error {
	table, err := db.FindTable(params.ctx, params.tableID)
	if err != nil {
		log.Errorf("Couldn't reserve seat for %s, finding table: %s", params.iden, err)
		return err
	}
	err = table.ReserveSeat(params.position, params.iden)
	if err != nil {
		return err
	}

	seat := table.GetSeatUnsafe(params.position)

	err = db.SetTableReservation(params.ctx, table, seat)
	if err != nil {
		if !errors.Is(err, db.ErrVersionNotMatch) {
			log.Warnf("reserve seat race condition %s: %s", params.iden, err)
		}
		return err
	}

	broadcast.SendTableEvent(params.tableID.Hex(), events.BuildSeatReservedEvent(seat))
	actionhandler.LaunchSeatReservationTimeout(timeout.Key{TableID: table.ID, Position: params.position, Version: seat.Version}, seat.ReservationTimeoutAt())

	return nil
}

func CancelSeatReservation(params *PositionParams) error {
	ctx := params.ctx
	tableID := params.tableID.Hex()
	iden := params.iden
	position := params.position

	table, err := db.FindTable(ctx, params.tableID)
	if err != nil {
		log.Errorf("Couldn't reserve seat for %s, finding table: %s", iden, err)
		return err
	}

	seat, err := table.GetSeat(position)
	if err != nil {
		return services.ErrFormat("not seat at position")
	}
	playerToRemove := seat.GetPlayer()

	err = table.CancelSeatReservation(position, iden)
	if err != nil {
		return err
	}

	err = db.CancelTableReservation(ctx, table, seat)
	if err != nil {
		return err
	}

	broadcast.SendTableEvent(tableID, events.BuildPlayerLeft(playerToRemove))

	return nil
}
