package db

import (
	"context"
	conf "github.com/glossd/pokergloss/goconf"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const DefaultTimeout = time.Second
const DbName = "bank"

var Client *mongo.Client

func Init() (context.Context, *mongo.Client, error) {
	return InitWithURI(conf.GetDbURI(DbName), true)
}

func InitWithURI(uri string, migrations bool) (context.Context, *mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return ctx, nil, err
	}

	if migrations {
		runMigrations(uri)
	}

	log.Infof("bank connected to MongoDB on %s", conf.Props.DB.Host)
	return ctx, Client, nil
}
