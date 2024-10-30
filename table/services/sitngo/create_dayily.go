package sitngo

import (
	"context"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var simpleID1 = oid("61dc108d8c8e870ed4f1613e")
var simpleID2 = oid("61dc1097d2c09fc218cdbb37")
var mediumID = oid("61dc10aa92402c73317ee178")
var PersistentSitNGoIDs = map[primitive.ObjectID]struct{}{
	simpleID1: {},
	simpleID2: {},
	mediumID:  {},
}

func CreateDaily() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var lobbiesToInsert []*domain.LobbySitAndGo

	_, err := db.FindSitAndGoLobby(ctx, simpleID1)
	if err == mongo.ErrNoDocuments {
		lobbiesToInsert = append(lobbiesToInsert, lobbySitNGo(simpleID1, "Sit&Go #1", 1000))
	}
	_, err = db.FindSitAndGoLobby(ctx, simpleID2)
	if err == mongo.ErrNoDocuments {
		lobbiesToInsert = append(lobbiesToInsert, lobbySitNGo(simpleID2, "Sit&Go #2", 1000))
	}

	_, err = db.FindSitAndGoLobby(ctx, mediumID)
	if err == mongo.ErrNoDocuments {
		lobbiesToInsert = append(lobbiesToInsert, lobbySitNGo(mediumID, "Medium #1", 5000))
	}

	if len(lobbiesToInsert) == 0 {
		return nil
	}

	return db.InsertManySitNGoLobbies(ctx, lobbiesToInsert)
}

func lobbySitNGo(id primitive.ObjectID, name string, buyIn int64) *domain.LobbySitAndGo {
	l, err := domain.NewLobbySitAndGo(domain.NewLobbySitAndGoParams{
		NewTableParams: domain.NewTableParams{
			Name:            name,
			Size:            3,
			BigBlind:        10,
			DecisionTimeout: domain.DefaultDecisionTimeout,
			BettingLimit:    domain.NL,
			Identity:        domain.SystemIdentity,
		},
		PlacesPaid:        1,
		BuyIn:             buyIn,
		LevelIncreaseTime: domain.DefaultLevelIncreaseTime,
	})
	if err != nil {
		log.Fatalf("Failed to create new sitngo lobby: %s", err)
	}
	l.ID = id
	return l
}

func oid(s string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		log.Fatal(err)
	}
	return oid
}
