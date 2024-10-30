package goconf

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDbURI(dbName string) string {
	dbConf := Props.DB
	fullHost := dbConf.Host
	if dbConf.Port != nil {
		fullHost = fmt.Sprintf("%s:%d", fullHost, *dbConf.Port)
	}

	var creds string
	if dbConf.Username != "" {
		creds = fmt.Sprintf("%s:%s@", dbConf.Username, dbConf.Password)
	}

	return fmt.Sprintf("%s://%s%s/%s?retryWrites=true&w=majority",
		dbConf.Scheme,
		creds,
		fullHost,
		dbName)
}

func MustObjectIDFromHex(s string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return id
}
