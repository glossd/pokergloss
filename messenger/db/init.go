package db

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const DefaultTimeout = time.Second

var Client *mongo.Client

const DbName = "messenger"

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

func conf.GetDbURI(DbName) string {
	dbConf := conf.Props.DB
	fullHost := dbConf.Host
	if dbConf.Port != nil {
		fullHost = fmt.Sprintf("%s:%d", fullHost, *dbConf.Port)
	}

	var creds string
	if dbConf.Username != "" {
		creds = fmt.Sprintf("%s:%s", dbConf.Username, dbConf.Password)
	}

	return fmt.Sprintf("%s://%s@%s/%s?retryWrites=true&w=majority",
		dbConf.Scheme,
		creds,
		fullHost,
		DbName)
}

func filterID(id interface{}) bson.M {
	return bson.M{"_id": id}
}
