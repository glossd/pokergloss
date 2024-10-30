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

var Client *mongo.Client

const DbName = "profile"

func Init() (context.Context, *mongo.Client, error) {
	return InitWithURI(conf.GetDbURI(DbName))
}

func InitWithURI(uri string) (context.Context, *mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return ctx, nil, err
	}
	log.Infof("Connected to MongoDB on %s", conf.Props.DB.Host)

	runMigrations(uri)

	return ctx, Client, nil
}
