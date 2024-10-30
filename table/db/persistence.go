package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func PersistenceCol() *mongo.Collection {
	return Client.Database(DbName).Collection("persistence")
}
