package mq

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

// Can be deleted if there'll be enough money for good mongo.
// Problem came with cheap database. Only for DecisionTimeout. This map stores existed gameflow version. Prevents expired DecisionTimeout
// to make a database request. Could scale horizontally, but won't be as effective.
var gameFlowVersionCache = sync.Map{} // tableID -> gameFlowVersion

func SetCacheGameFlow(tableID primitive.ObjectID, version int64) {
	gameFlowVersionCache.Store(tableID.Hex(), version)
}

func GetCacheGameFlow(tableID primitive.ObjectID) int64 {
	version, ok := gameFlowVersionCache.Load(tableID.Hex())
	if !ok {
		return 0
	}
	return version.(int64)
}
