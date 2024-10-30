package service

import (
	"context"
	"github.com/glossd/pokergloss/achievement/db"
	"github.com/glossd/pokergloss/achievement/domain"
	"github.com/glossd/pokergloss/achievement/web/mq/bank"
	"github.com/glossd/pokergloss/achievement/web/mq/ws"
	"github.com/glossd/pokergloss/gomq/mqtable"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserExp(ctx context.Context, userID string) (*domain.ExP, error) {
	exp, err := db.FindExP(ctx, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.NewExP(userID), nil
		}
		return nil, err
	}

	return exp, nil
}

func UpdateExp(gameEnd *mqtable.GameEnd) {
	for _, player := range gameEnd.Players {
		exp, err := db.FindExpNoCtx(player.UserId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				exp = domain.NewExP(player.UserId)
			} else {
				log.Errorf("ExP update error, player=%+v: %s", player, err)
				return
			}
		}

		exp.UpdateWithGameEnd(gameEnd)
		err = db.UpsertExpNoCtx(exp)
		if err != nil {
			return
		}

		ws.PublishExpEvent(exp)
		bank.DepositNewLevel(exp)
	}
}
