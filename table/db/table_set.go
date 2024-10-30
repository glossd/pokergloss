package db

import (
	"context"
	"github.com/glossd/pokergloss/table/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetTableOneContext(ctx context.Context, tableID primitive.ObjectID, element bson.E) error {
	update := bson.D{{"$set", bson.D{element}}}
	return updateTableContext(ctx, tableID, update)
}

func SetTableOne(tableID primitive.ObjectID, element bson.E) error {
	update := bson.D{{"$set", bson.D{element}}}
	return updateTable(tableID, update)
}

func SetTableContext(ctx context.Context, tableID primitive.ObjectID, elements []bson.E) error {
	update := bson.D{{"$set", elements}}
	return updateTableContext(ctx, tableID, update)
}

func SetTable(tableID primitive.ObjectID, elements []bson.E) error {
	update := bson.D{{"$set", elements}}
	return updateTable(tableID, update)
}

func SetTableGameFlow(ctx context.Context, tableID primitive.ObjectID, version int64, elements []bson.E) error {
	update := bson.D{{"$set", elements}}
	filter := bson.D{ID(tableID), {"gameflowversion", version}}
	res, err := ColTable().UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("Update table, filter=%v failed: %v", filter, err)
		return err
	}
	if res.MatchedCount == 0 {
		log.Warnf("Race condition for game flow happened, tableID=%s", tableID.Hex())
		return ErrVersionNotMatch
	}
	return nil
}

func SetTableGameFlowNoCtx(tableID primitive.ObjectID, version int64, elements []bson.E) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return SetTableGameFlow(ctx, tableID, version, elements)
}

func GameFlowVersion(table *domain.Table) bson.E {
	return bson.E{Key: "gameflowversion", Value: table.GameFlowVersion}
}

func SetTableUpdates(ctx context.Context, tableID primitive.ObjectID, elements ...bson.E) error {
	return updateTableContext(ctx, tableID, elements)
}

func SetTableFilter(ctx context.Context, tableID primitive.ObjectID, filter []bson.E, updates []bson.E) error {
	filter = append(filterID(tableID), filter...)
	update := bson.D{{"$set", updates}}
	result, err := ColTable().UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("Update table %s filter failed: %s", tableID.Hex(), err)
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func updateTable(tableID primitive.ObjectID, update bson.D) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return updateTableContext(ctx, tableID, update)
}

func updateTableContext(ctx context.Context, tableID primitive.ObjectID, update bson.D) error {
	_, err := ColTable().UpdateOne(ctx, filterID(tableID), update)
	if err != nil {
		log.Errorf("Update table %s failed: %s", tableID.Hex(), err)
		return err
	}
	return nil
}
