package db

import (
	"context"
	"github.com/glossd/pokergloss/table/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func FindSitAndGoLobby(ctx context.Context, id primitive.ObjectID) (*domain.LobbySitAndGo, error) {
	var result domain.LobbySitAndGo
	err := ColSitAndGoLobby().FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		log.Errorf("Find sit&go lobby failed: %s", err)
		return nil, err
	}
	return &result, nil
}

func FindSitAndGoLobbiesForStartAt(ctx context.Context, startAt int64) ([]*domain.LobbySitAndGo, error) {
	cur, err := ColSitAndGoLobby().Find(ctx, bson.D{{"startat", startAt}, {"status", domain.LobbyRegistering}})
	if err != nil {
		log.Errorf("Find sit&go lobbies by startAt failed: %s", err)
		return nil, err
	}
	var res []*domain.LobbySitAndGo
	for cur.Next(ctx) {
		var l domain.LobbySitAndGo
		err := cur.Decode(&l)
		if err != nil {
			log.Errorf("lobby sitngo decode failed: %s", err)
			return nil, err
		}
		res = append(res, &l)
	}
	cur.Close(ctx)
	return res, nil
}

func FindSitAndGoLobbyNoCtx(id primitive.ObjectID) (*domain.LobbySitAndGo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return FindSitAndGoLobby(ctx, id)
}

func InsertSitAndGoLobby(ctx context.Context, l *domain.LobbySitAndGo) error {
	l.PlayersCount = len(l.Entries)
	_, err := ColSitAndGoLobby().InsertOne(ctx, l)
	if err != nil {
		log.Errorf("Couldn't insert sit&go lobby: %s", err)
		return err
	}
	return nil
}

func InsertSitAndGoLobbyNoCtx(l *domain.LobbySitAndGo) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	l.PlayersCount = len(l.Entries)
	_, err := ColSitAndGoLobby().InsertOne(ctx, l)
	if err != nil {
		log.Errorf("Couldn't insert sit&go lobby: %s", err)
		return err
	}
	return nil
}

func InsertManySitNGoLobbies(ctx context.Context, tables []*domain.LobbySitAndGo) error {
	adapted := make([]interface{}, 0, len(tables))
	for _, table := range tables {
		adapted = append(adapted, table)
	}

	_, err := ColSitAndGoLobby().InsertMany(ctx, adapted)
	if err != nil {
		log.Errorf("Failed insert many sitngo lobbies: %s", err)
		return err
	}
	return nil
}

func UpdateSitngoLobbyStatus(id primitive.ObjectID, status domain.LobbyStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := ColSitAndGoLobby().UpdateOne(ctx, filterID(id), bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		log.Errorf("Failed set sitngo lobby to finished, id=%s", id.Hex())
		return err
	}
	return nil
}

func UpdateSitAndGoLobby(ctx context.Context, l *domain.LobbySitAndGo) error {
	l.PlayersCount = len(l.Entries)
	filter := bson.D{{"_id", l.ID}, {"version", l.Version}}
	l.Version++
	res, err := ColSitAndGoLobby().ReplaceOne(ctx, filter, l)
	if err != nil {
		log.Errorf("Update sit&go lobby %s failed: %s", l.ID.Hex(), err)
		return err
	}
	if res.MatchedCount == 0 {
		return ErrVersionNotMatch
	}
	return nil
}

func ForEachSitngoLobby(filter interface{}, apply func(t *domain.LobbySitAndGo)) error {
	cur, err := ColSitAndGoLobby().Find(context.Background(), filter)
	if err != nil {
		log.Errorf("Find tables failed: %s", err)
		return err
	}

	curCtx := context.Background()
	for cur.Next(curCtx) {
		var l domain.LobbySitAndGo
		err := cur.Decode(&l)
		if err != nil {
			log.Errorf("lobby sitngo decode failed: %s", err)
			return err
		}
		apply(&l)
	}
	cur.Close(curCtx)
	return nil
}

func DeleteSitngo(ctx context.Context, id primitive.ObjectID) error {
	_, err := ColSitAndGoLobby().DeleteOne(ctx, filterID(id))
	if err != nil {
		log.Errorf("Failed to delete sitngo lobby %s : %s", id.Hex(), err)
		return err
	}
	return nil
}

func FilterSitngoLobbies(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]*domain.LobbySitAndGo, error) {
	// A slice of tables for storing the decoded documents
	var results []*domain.LobbySitAndGo

	cur, err := ColSitAndGoLobby().Find(ctx, filter, opts...)
	if err != nil {
		log.Errorf("Filter sit&go lobby failed: %s", err)
		return nil, err
	}

	for cur.Next(ctx) {
		var r domain.LobbySitAndGo
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
		return []*domain.LobbySitAndGo{}, nil
	}

	return results, nil
}

func ColSitAndGoLobby() *mongo.Collection {
	return Client.Database(DbName).Collection("sitAndGoLobbies")
}
