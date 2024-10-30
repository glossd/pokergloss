package tables

import (
	"context"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services"
	"github.com/glossd/pokergloss/table/services/paging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindAll(ctx context.Context, tableType domain.TableType, params paging.Params) ([]*domain.Table, error) {
	var filter []bson.E
	filter = append(filter, bson.E{Key: "type", Value: tableType})
	filter = append(filter, bson.E{Key: "isprivate", Value: false})
	filter = append(filter, db.SkipEmptyFullFilter(params)...)
	opts := db.PagingOptions(params)
	opts.SetSort(bson.D{{"_id", 1}})
	return db.FilterTables(ctx, filter, opts)
}

func Find(ctx context.Context, tableID string) (*domain.Table, error) {
	oid, err := primitive.ObjectIDFromHex(tableID)
	if err != nil {
		return nil, services.E(err)
	}
	table, err := db.FindTable(ctx, oid)
	if err != nil {
		return nil, err
	}
	return table, nil
}

func Create(ctx context.Context, params domain.NewTableParams) (*domain.Table, error) {
	t, err := domain.NewTable(params)
	if err != nil {
		return nil, err
	}
	err = db.InsertTable(ctx, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func Delete(ctx context.Context, ID string) error {
	oid, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return services.E(err)
	}
	err = db.DeleteTable(ctx, oid)
	if err != nil {
		return err
	}
	return nil
}
