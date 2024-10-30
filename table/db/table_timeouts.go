package db

import (
	"context"
	"github.com/glossd/pokergloss/table/domain"
	"go.mongodb.org/mongo-driver/bson"
)

func FindTablesUnfinishedTimeout() ([]*domain.Table, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return FilterTables(ctx, bson.D{
		{"$or", bson.A{
			bson.D{{"status", bson.D{{"$ne", domain.WaitingTable}}}},
			bson.D{{"seats.player.status", domain.PlayerReservedSeat}},
		}},
	})

}
