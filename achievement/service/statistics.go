package service

import (
	"context"
	"github.com/glossd/pokergloss/achievement/db"
	"github.com/glossd/pokergloss/achievement/domain"
	"github.com/glossd/pokergloss/gomq/mqtable"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserStatistics(ctx context.Context, userID string) (*domain.Statistics, error) {
	stat, err := db.FindStatistics(ctx, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.NewStatistics(userID), nil
		}
		return nil, err
	}

	return stat, nil
}

func UpdateStatistics(ge *mqtable.GameEnd) {
	gameEnd := domain.NewGameEnd(ge)
	for _, player := range ge.Players {
		stat, err := db.FindStatisticsNoCtx(player.UserId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				stat = domain.NewStatistics(player.UserId)
			} else {
				log.Errorf("Failed to find statistics: %s", err)
				return
			}
		}
		stat.Update(gameEnd)
		err = db.UpsertStatistics(stat)
		if err != nil {
			log.Errorf("Failed to find statistics: %s", err)
			return
		}
	}
}
