package e2e

import (
	"context"
	"github.com/glossd/pokergloss/table-history/db"
	"github.com/pokerblow/mongotest"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	cc := mongotest.StartMongoContainer("4.2")

	_, _, err := db.InitWithURI(cc.GetMongoURI(db.DbName))
	if err != nil {
		cc.KillMongoContainer()
		log.Fatal(err)
	}

	code := m.Run()

	cc.KillMongoContainer()

	os.Exit(code)
}

func cleanUpDB() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := db.ColTables().Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
