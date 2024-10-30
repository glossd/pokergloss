package db

import (
	"context"
	"github.com/glossd/pokergloss/assignment/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindUserDaily(ctx context.Context, userID string) (*domain.UserDaily, error) {
	var d domain.UserDaily
	err := ColUserDaily().FindOne(ctx, idFilter(userID)).Decode(&d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func UpsertUserDaily(ctx context.Context, ud *domain.UserDaily) (err error) {
	if ud.Version == 0 {
		_, err = ColUserDaily().ReplaceOne(ctx, idFilter(ud.UserID), ud, &options.ReplaceOptions{Upsert: &True})
	} else {
		ud.Version++
		_, err = ColUserDaily().ReplaceOne(ctx, idVersionFilter(ud.UserID, ud.Version-1), ud)
	}
	if err != nil {
		log.Errorf("db.UpsertUserDaily failed: %s", err)
		return err
	}
	return nil
}

func ColUserDaily() *mongo.Collection {
	return Client.Database(DbName).Collection("userDaily")
}
