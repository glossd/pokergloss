package e2e

import (
	"context"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/web/client/mq"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func cleanUp() {
	cleanUpDb()
	mq.ResetTestMQ()
}

func cleanUpDb() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := db.ColTable().Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = db.ColSitAndGoLobby().Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = db.ColLobbyMulti().Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = db.ColLobbyMultiConfig().Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertTable(t *testing.T) *domain.Table {
	return InsertTableTimeout(t, -1)
}

func InsertTableTimeout(t *testing.T, timeout time.Duration) *domain.Table {
	table := NewTableTimeout(t, timeout)
	insertTable(t, table)
	return table
}

func NewTable(t *testing.T) *domain.Table {
	return NewTableTimeout(t, -1)
}

func NewTableTimeout(t *testing.T, timeout time.Duration) *domain.Table {
	params := tableParams(9, timeout)
	table, err := domain.NewTable(params)
	assert.Nil(t, err)
	return table
}

func tableParams(size int, timeout time.Duration) domain.NewTableParams {
	return domain.NewTableParams{
		Name:            "my table",
		Size:            size,
		BigBlind:        defaultBigBlind,
		DecisionTimeout: timeout,
		Identity:        defaultIdentity,
	}
}

func insertTable(t *testing.T, table *domain.Table) {
	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := db.InsertTable(timeout, table)
	assert.Nil(t, err)
}
