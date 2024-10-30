package db

import (
	"context"
	"github.com/glossd/pokergloss/bank/db/paging"
	"github.com/glossd/pokergloss/bank/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func InsertOperation(ctx context.Context, operation *domain.Operation) error {
	operation.ID = primitive.NewObjectID()
	_, err := OperationCol().InsertOne(ctx, operation)
	if err != nil {
		return err
	}
	return nil
}

func FindOperationNoCtx(id primitive.ObjectID) (*domain.Operation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return FindOperation(ctx, id)
}

func FindOperation(ctx context.Context, id primitive.ObjectID) (*domain.Operation, error) {
	operations, err := filterOperations(ctx, filterID(id))
	if err != nil {
		return nil, err
	}
	if len(operations) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return operations[0], nil
}

func FindOperations(ctx context.Context) ([]*domain.Operation, error) {
	return filterOperations(ctx, bson.D{{}})
}

func FindAllOperationsByUserId(userID string) ([]*domain.Operation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return filterOperations(ctx, bson.D{{"userid", userID}})
}

func FindOperationsByUserIdReverse(ctx context.Context, userID string, page paging.Page) ([]*domain.Operation, error) {
	return filterOperations(ctx,
		bson.D{{"userid", userID}},
		&options.FindOptions{Limit: &page.Limit, Skip: &page.Skip, Sort: bson.M{"_id": -1}})
}

func filterOperations(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]*domain.Operation, error) {
	// A slice of tables for storing the decoded documents
	var tables []*domain.Operation

	cur, err := OperationCol().Find(ctx, filter, opts...)
	if err != nil {
		return tables, err
	}

	for cur.Next(ctx) {
		var t domain.Operation
		err := cur.Decode(&t)
		if err != nil {
			return tables, err
		}
		tables = append(tables, &t)
	}

	if err := cur.Err(); err != nil {
		return tables, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(tables) == 0 {
		return []*domain.Operation{}, nil
	}

	return tables, nil
}

func OperationCol() *mongo.Collection {
	return Client.Database(DbName).Collection("operations")
}
