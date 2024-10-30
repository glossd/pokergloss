package db

import (
	"context"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/paging"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// usage examples of golang mongodb driver
// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
// https://www.digitalocean.com/community/tutorials/how-to-use-go-with-mongodb-using-the-mongodb-go-driver

const collectionName = "tables"

func FindTable(ctx context.Context, id primitive.ObjectID) (*domain.Table, error) {
	tables, err := FilterTables(ctx, bson.D{{"_id", id}})
	if err != nil {
		return nil, err
	}
	if len(tables) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return tables[0], nil
}

func FindTableNoCtx(id primitive.ObjectID) (*domain.Table, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return FindTable(ctx, id)
}

func FindTableGameFlowNoCtx(id primitive.ObjectID, gameFlowVersion int64) (*domain.Table, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return FindTableGameFlow(ctx, id, gameFlowVersion)
}

func FindTableGameFlow(ctx context.Context, id primitive.ObjectID, gameFlowVersion int64) (*domain.Table, error) {
	var result domain.Table
	err := ColTable().FindOne(ctx, bson.D{ID(id), {Key: "gameflowversion", Value: gameFlowVersion}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrVersionNotMatch
		}
		log.Errorf("Failed to fetch table game flow, tableID=%s: %s", id.Hex(), err)
		return nil, err
	}
	return &result, nil
}

func FindTablesByLobbyID(lobbyID primitive.ObjectID) ([]*domain.Table, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return FilterTables(ctx, bson.M{"tournamentattributes.lobbyid": lobbyID})
}

func ForEachTable(filter interface{}, apply func(ctx context.Context, t *domain.Table)) error {
	cur, err := ColTable().Find(context.Background(), filter)
	if err != nil {
		log.Errorf("Find tables failed: %s", err)
		return err
	}

	curCtx := context.Background()
	for cur.Next(curCtx) {
		var t domain.Table
		err := cur.Decode(&t)
		if err != nil {
			log.Errorf("document iteration failed: %s", err)
			return err
		}
		apply(curCtx, &t)
	}
	cur.Close(curCtx)
	return nil
}

func FindTablesContext(ctx context.Context, tableType domain.TableType, params paging.Params) ([]*domain.Table, error) {
	var filter []bson.E
	filter = append(filter, bson.E{Key: "type", Value: tableType})
	filter = append(filter, SkipEmptyFullFilter(params)...)
	return FilterTables(ctx, filter, PagingOptions(params))
}

func FilterTables(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]*domain.Table, error) {
	// A slice of tables for storing the decoded documents
	var tables []*domain.Table

	cur, err := ColTable().Find(ctx, filter, opts...)
	if err != nil {
		log.Errorf("Filter tables failed: %s", err)
		return nil, err
	}

	for cur.Next(ctx) {
		var t domain.Table
		err := cur.Decode(&t)
		if err != nil {
			return tables, err
		}
		tables = append(tables, &t)
	}

	if err := cur.Err(); err != nil {
		log.Errorf("Filter tables failed: %s", err)
		return tables, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(tables) == 0 {
		return []*domain.Table{}, nil
	}

	return tables, nil
}

func ColTable() *mongo.Collection {
	return Client.Database(DbName).Collection(collectionName)
}
