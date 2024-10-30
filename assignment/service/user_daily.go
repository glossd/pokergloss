package service

import (
	"context"
	"github.com/glossd/pokergloss/assignment/db"
	"github.com/glossd/pokergloss/assignment/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUserDaily(ctx context.Context, userID string) (*domain.UserDaily, error) {
	d, err := db.FindDailyOfNow(ctx)
	if err != nil {
		return nil, err
	}
	return FindOrBuildUserDaily(ctx, userID, d)
}

func FindOrBuildUserDaily(ctx context.Context, userID string, d *domain.Daily) (*domain.UserDaily, error) {
	ud, err := db.FindUserDaily(ctx, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.NewUserDaily(userID, d), nil
		}
		log.Errorf("db.FindUserDaily failed: %s", err)
		return nil, err
	}
	if ud.DailyID != d.Day {
		return domain.NewUserDaily(userID, d), nil
	}
	return ud, nil
}
