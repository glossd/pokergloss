package db

import (
	"context"
	"github.com/glossd/pokergloss/assignment/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var daily *domain.Daily

func FindDailyOfNow(ctx context.Context) (*domain.Daily, error) {
	return FindDaily(ctx, domain.NewDailyID(time.Now()))
}

func FindDaily(ctx context.Context, id domain.DailyID) (*domain.Daily, error) {
	if daily != nil && daily.Day == id {
		return daily, nil
	}
	var d domain.Daily
	err := ColDaily().FindOne(ctx, idFilter(id)).Decode(&d)
	if err != nil {
		return nil, err
	}
	daily = &d
	return &d, nil
}

func InsertDaily(d *domain.Daily) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := ColDaily().InsertOne(ctx, d)
	if err != nil {
		return err
	}
	return nil
}

func ColDaily() *mongo.Collection {
	return Client.Database(DbName).Collection("daily")
}

func CleanCache() {
	daily = nil
}
