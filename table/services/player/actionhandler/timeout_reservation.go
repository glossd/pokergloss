package actionhandler

import (
	"context"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/services/broadcast"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	"github.com/glossd/pokergloss/table/web/client/mqpub"
	log "github.com/sirupsen/logrus"
)

func LaunchSeatReservationTimeout(key timeout.Key, timeoutAt int64) {
	if conf.Props.SeatReservationTimeout <= 0 {
		// for tests do it yourself
		return
	}
	mqpub.PublishTimeoutEvent(&timeout.Event{
		Type: timeout.SeatReservation,
		At:   timeoutAt,
		Key:  key,
	})
}

func DoSeatReservationTimeout(ctx context.Context, key timeout.Key) (tryAgain bool) {
	table, err := db.FindTable(ctx, key.TableID)
	if err != nil {
		log.Errorf("Couldn't cancel reservation, finding tableID=%s, position=%d : %s", key.TableID, key.Position, err)
		return true
	}

	seat, err := table.GetSeat(key.Position)
	if err != nil {
		log.Tracef("Reservation timeout of empty seat tableID=%s, position=%d : %s", key.TableID, key.Position, err)
		return
	}

	if seat.Version != key.Version {
		log.Tracef("Reservation timeout race condition tableID=%s, position=%d : %s", key.TableID, key.Position, err)
		return
	}

	err = table.CancelSeatReservationTimeout(key.Position)
	if err != nil {
		log.Errorf("Couldn't cancel reservation tableID=%s, position=%d : %s", key.TableID, key.Position, err)
		return
	}

	err = db.CancelTableReservation(ctx, table, seat)
	if err != nil {
		if err == db.ErrVersionNotMatch {
			log.Debugf("Seat has changed before timeout, key=%s", key)
			return
		}
		log.Errorf("Couldn't cancel reservation tableID=%s, position=%d : %s", key.TableID, key.Position, err)
		return true
	}

	broadcast.SendTableEvent(key.TableID.Hex(), events.BuildSeatReservationTimeout(key.Position))
	return
}
