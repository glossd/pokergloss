package service

import (
	"context"
	"github.com/glossd/pokergloss/assignment/db"
	"github.com/glossd/pokergloss/assignment/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func CreateDaily(t time.Time) (*domain.Daily, error) {
	daily := domain.NewDaily(t)
	err := db.InsertDaily(daily)
	if err != nil {
		log.Errorf("CreateDaily failed: %s", err)
		return nil, err
	}
	return daily, nil
}

func RecoverDaily() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := db.FindDaily(ctx, domain.NewDailyID(time.Now()))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err := CreateDaily(time.Now())
			if err != nil {
				time.AfterFunc(5*time.Second, RecoverDaily)
			} else {
				log.Infof("Successfully recovered daily")
			}
			return
		} else {
			log.Errorf("FindDaily failed: %s", err)
			time.AfterFunc(5*time.Second, RecoverDaily)
		}
	}
}
