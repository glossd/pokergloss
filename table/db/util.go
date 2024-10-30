package db

import (
	"errors"
	"fmt"
	"github.com/glossd/pokergloss/table/services/paging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrVersionNotMatch = errors.New("version of the resource doesn't match")

func SeatDbPath(position int) string {
	return fmt.Sprintf("seats.%d", position)
}

func PlayerDbPath(position int) string {
	return fmt.Sprintf("seats.%d.player", position)
}

func filterID(id interface{}) bson.D {
	return bson.D{{"_id", id}}
}

func ID(id primitive.ObjectID) bson.E {
	return bson.E{Key: "_id", Value: id}
}

func SkipEmptyFullFilter(params paging.Params) bson.D {
	var filterAcc bson.A
	if params.SkipEmpty {
		filterAcc = append(filterAcc, bson.D{
			{Key: "$expr", Value: bson.M{"$ne": bson.A{0, "$playerscount"}}},
		})
	}
	if params.SkipFull {
		filterAcc = append(filterAcc, bson.D{
			{Key: "$expr", Value: bson.M{"$ne": bson.A{"$size", "$playerscount"}}},
		})
	}

	filter := bson.D{}
	if len(filterAcc) > 0 {
		filter = bson.D{{"$and", filterAcc}}
	}
	return filter
}

func PagingOptions(params paging.Params) *options.FindOptions {
	return &options.FindOptions{Skip: &params.Skip, Limit: &params.Limit}
}
