package db

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type RebalancerConfig struct {
	LobbyID     primitive.ObjectID `bson:"_id"`
	CountDown   int
	LockOnMsgId string
}

func UpsertRebalanceConfig(ctx context.Context, rc *RebalancerConfig) error {
	_, err := RebalanceConfigCol().ReplaceOne(ctx, filterID(rc.LobbyID), rc, &options.ReplaceOptions{Upsert: &True})
	if err != nil {
		log.Errorf("Failed to insert RebalancerConfig: %s", err)
		return err
	}
	return nil
}

func UpdateRebalanceConfig(ctx context.Context, rc *RebalancerConfig) error {
	_, err := RebalanceConfigCol().ReplaceOne(ctx, filterID(rc.LobbyID), rc)
	if err != nil {
		log.Errorf("Failed to update RebalancerConfig: %s", err)
		return err
	}
	return nil
}

func FindRebalanceConfigAndCountDown(ctx context.Context, msgID string, lobbyID primitive.ObjectID) (*RebalancerConfig, error) {
	returnDoc := options.After
	opts := options.FindOneAndUpdateOptions{ReturnDocument: &returnDoc}
	var rc RebalancerConfig
	orFilter := bson.E{Key: "$or", Value: bson.A{bson.M{"lockonmsgid": ""}, bson.M{"lockonmsgid": msgID}}}
	err := RebalanceConfigCol().FindOneAndUpdate(ctx, bson.D{{"_id", lobbyID}, orFilter}, bson.M{"$inc": bson.M{"countdown": -1}, "$set": bson.M{"lockonmsgid": msgID}}, &opts).Decode(&rc)
	if err != nil {
		return nil, err
	}

	return &rc, nil
}

func RecursiveUnlockRebalanceConfig(lobbyID primitive.ObjectID) {
	err := UnlockRebalanceConfig(lobbyID)
	if err != nil && err != mongo.ErrNoDocuments {
		time.Sleep(time.Second)
		RecursiveUnlockRebalanceConfig(lobbyID)
	}
}

func UnlockRebalanceConfig(lobbyID primitive.ObjectID) error {
	_, err := RebalanceConfigCol().UpdateOne(context.Background(), filterID(lobbyID), bson.M{"$set": bson.M{"lockonmsgid": ""}})
	if err != nil {
		log.Errorf("Failed to unlock rebalance config: %s", err)
		return err
	}
	return nil
}

func RebalanceConfigCol() *mongo.Collection {
	return Client.Database(DbName).Collection("rebalanceConfig")
}
