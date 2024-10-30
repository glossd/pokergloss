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

func InsertManyLobbyMulti(lobbies []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := ColLobbyMulti().InsertMany(ctx, lobbies)
	if err != nil {
		return err
	}
	return nil
}

func InsertOneLobbyMulti(lobbies *domain.LobbyMulti) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := ColLobbyMulti().InsertOne(ctx, lobbies)
	if err != nil {
		return err
	}
	return nil
}

func FindMultiLobbiesSortedByStartAt(ctx context.Context, params paging.Params) ([]*domain.LobbyMulti, error) {
	findOptions := PagingOptions(params)
	findOptions.SetSort(bson.D{{"startat", 1}, {"_id", 1}})
	return filterMultiLobbies(ctx, SkipEmptyFullFilter(params), findOptions)
}

func FindMultiLobbies(ctx context.Context, params paging.Params) ([]*domain.LobbyMulti, error) {
	return filterMultiLobbies(ctx, SkipEmptyFullFilter(params), PagingOptions(params))
}

func FindMultiLobbiesNoCtx(params paging.Params) ([]*domain.LobbyMulti, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return FindMultiLobbies(ctx, params)
}

func FindAllMultiLobbiesNoCtx() ([]*domain.LobbyMulti, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return filterMultiLobbies(ctx, bson.D{})
}

func FindLobbyMulti(ctx context.Context, id primitive.ObjectID) (*domain.LobbyMulti, error) {
	var result domain.LobbyMulti
	err := ColLobbyMulti().FindOne(ctx, filterID(id)).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func FindLobbyMultiNoCtx(id primitive.ObjectID) (*domain.LobbyMulti, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return FindLobbyMulti(ctx, id)
}

func FindFirstLobbyMultiNoCtx() (*domain.LobbyMulti, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	var result domain.LobbyMulti
	err := ColLobbyMulti().FindOne(ctx, bson.D{}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func UpdateLobbyMultiStatusNoCtx(id primitive.ObjectID, status domain.LobbyStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	_, err := ColLobbyMulti().UpdateOne(ctx, filterID(id), bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		log.Errorf("Failed to update multi lobby status: %s", err)
	}
	return err
}

func GetLobbyMultiByStartAt(startAt int64) (*domain.LobbyMulti, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	var lobby domain.LobbyMulti
	err := ColLobbyMulti().FindOne(ctx, bson.M{"startat": startAt}).Decode(&lobby)
	if err != nil {
		return nil, err
	}

	return &lobby, nil
}

func FindMultiLobbiesByStartAt(ctx context.Context, startAt int64) ([]*domain.LobbyMulti, error) {
	return filterMultiLobbies(ctx, bson.M{"startat": startAt})
}

func UpdateLobbyMulti(ctx context.Context, l *domain.LobbyMulti) error {
	filter := bson.D{{"_id", l.ID}, {"version", l.Version}}
	l.Version++
	res, err := ColLobbyMulti().ReplaceOne(ctx, filter, l)
	if err != nil {
		log.Errorf("Update multi lobby %s failed: %s", l.ID.Hex(), err)
		return err
	}
	if res.MatchedCount == 0 {
		return ErrVersionNotMatch
	}
	return nil
}

func UpdateLobbyMultiNoCtx(l *domain.LobbyMulti) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return UpdateLobbyMulti(ctx, l)
}

func DeleteLobbyMulti(ctx context.Context, id primitive.ObjectID) error {
	_, err := ColLobbyMulti().DeleteOne(ctx, filterID(id))
	if err != nil {
		log.Errorf("Failed to delete lobby multi, id=%s : %s", id.Hex(), err)
		return err
	}
	return nil
}

func filterMultiLobbies(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]*domain.LobbyMulti, error) {
	// A slice of tables for storing the decoded documents
	var results []*domain.LobbyMulti

	cur, err := ColLobbyMulti().Find(ctx, filter, opts...)
	if err != nil {
		log.Errorf("Filter sit&go lobby failed: %s", err)
		return nil, err
	}

	for cur.Next(ctx) {
		var r domain.LobbyMulti
		err := cur.Decode(&r)
		if err != nil {
			return results, err
		}
		results = append(results, &r)
	}

	if err := cur.Err(); err != nil {
		log.Errorf("Filter sit&go lobby failed: %s", err)
		return results, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(results) == 0 {
		return []*domain.LobbyMulti{}, nil
	}

	return results, nil
}

func ColLobbyMulti() *mongo.Collection {
	return Client.Database(DbName).Collection("multiLobbies")
}
