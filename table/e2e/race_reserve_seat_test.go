package e2e

import (
	"context"
	"github.com/glossd/pokergloss/auth/authid"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/player"
	"github.com/stretchr/testify/assert"
	"log"
	"sync/atomic"
	"testing"
	"time"
)

var satUsers = new(int32)

func TestReserveSeatRace(t *testing.T) {
	prevPropsSetup(t)
	conf.Props.Table.SeatReservationTimeout = time.Second

	table, err := domain.NewTable(domain.NewTableParams{
		Name:            "my table",
		Size:            9,
		BigBlind:        2,
		DecisionTimeout: 1,
		Identity:        defaultIdentity,
	})
	assert.Nil(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = db.InsertTable(ctx, table)
	assert.Nil(t, err)

	runWorkers(table, reserveSeat, concurrentUsers)
	assert.EqualValues(t, 1, *satUsers)
}

func reserveSeat(ctx context.Context, tableID string, userID string) {
	params, err := player.NewPositionParams(ctx, tableID, 0, authid.Identity{UserId: userID, Username: "username"})
	if err != nil {
		log.Fatal(err)
	}
	err = player.ReserveTableSeat(params)
	if err == nil {
		atomic.AddInt32(satUsers, 1)
	}
}
