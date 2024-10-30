package e2e

import (
	"github.com/glossd/pokergloss/auth"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/web/client/mqsub"
	"github.com/pokerblow/mongotest"
	"log"
	"os"
	"testing"
	"time"
)

const (
	defaultGameEndTimeout = 5 * time.Millisecond
)

func TestMain(m *testing.M) {
	conf.IsE2EVar = true
	conf.Props.Table.MinDecisionTimeout = -1
	conf.Props.Table.GameEndMinTimeout = 0
	conf.Props.SeatReservationTimeout = -1
	conf.Props.PlayerActionDuration = 0
	conf.Props.Enrich.PlayersEnabled = false
	conf.Props.RebalancerPeriod = time.Nanosecond
	conf.Props.Table.RakePercent = 0.0
	conf.Props.Tournament.FeePercent = 0.0
	domain.Algo = &domain.MockAlgo{}

	os.Setenv("PG_JWT_VERIFICATION_DISABLE", "true")
	auth.Init()

	go mqsub.SubscribeForTimeouts()
	go mqsub.SubscribeForMultiRebalance()

	cc := mongotest.StartMongoContainer("4.2")

	err := db.InitWithURI(cc.GetMongoURI(db.DbName))
	if err != nil {
		cc.KillMongoContainer()
		log.Fatal(err)
	}

	code := m.Run()

	cc.KillMongoContainer()

	os.Exit(code)
}
