package cash

import (
	"context"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type PersistenceCash struct {
	ID        string `bson:"_id"`
	CreatedAt int64
}

func CreatePersistentTables() {
	var configID = "cashPersistentTables"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var d PersistenceCash
	err := db.PersistenceCol().FindOne(ctx, bson.M{"_id": configID}).Decode(&d)
	if err == nil {
		return
	}
	if err != mongo.ErrNoDocuments {
		log.Errorf("Failed to create cash tables: %s", err)
		return
	}
	var tables = []*domain.Table{
		newPersistentTable("Mesopotamia", 2),
		newPersistentTable("Indus Valley", 10),
		newPersistentTable("Egypt", 40),
		newPersistentTable("Maya", 100),
		newPersistentTable("China", 400),
		newPersistentTable("Greece", 1000),
		newPersistentTable("Persia", 4000),
		newPersistentTable("Rome", 10000),
	}
	_, err = db.PersistenceCol().InsertOne(ctx, PersistenceCash{ID: configID, CreatedAt: time.Now().Unix()})
	if err != nil {
		log.Errorf("Failed to insert persistent tables, persistence: %s", err)
		return
	}
	err = db.InsertManyTables(ctx, tables)
	if err != nil {
		log.Errorf("Failed to insert persistent tables: %s", err)
		return
	}
}

func newPersistentTable(name string, bb int64) *domain.Table {
	table, err := domain.NewTable(domain.NewTableParams{
		Name:            name,
		Size:            6,
		BigBlind:        bb,
		DecisionTimeout: domain.DefaultDecisionTimeout,
		BettingLimit:    domain.NL,
		IsPrivate:       false,
		Identity:        domain.SystemIdentity,
	})
	if err != nil {
		log.Fatalf("cash.newTable %s", err)
	}

	table.IsPersistent = true
	return table
}
