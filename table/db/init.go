package db

import (
	"context"
	conf "github.com/glossd/pokergloss/goconf"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

const DefaultTimeout = time.Second

const DbName = "table"

var Client *mongo.Client

var lock sync.Mutex

func Init() {
	lock.Lock()
	defer lock.Unlock()
	if Client == nil {
		err := InitWithURI(conf.GetDbURI(DbName))
		if err != nil {
			log.Fatalf("Failed to init mongo client: %s", err)
		}
	}
}

func InitWithURI(uri string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	log.Infof("Connected to MongoDB on %s", conf.Props.DB.Host)

	runMigrations(uri)

	return nil
}
