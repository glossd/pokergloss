package daily

import (
	"context"
	"errors"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/bonus/db"
	"github.com/glossd/pokergloss/bonus/domain"
	"github.com/glossd/pokergloss/bonus/web/mq"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq/mqsurvival"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func TakeDaily(ctx context.Context, iden authid.Identity) (*domain.DailyBonus, error) {
	bonus, err := db.GetDailyBonus(ctx, iden.UserId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			bonus = domain.NewDailyBonus(iden.UserId)
		} else {
			return nil, err
		}
	}

	isAvailable := bonus.Visit()
	err = db.UpdateDailyBonus(ctx, bonus)
	if err != nil {
		return nil, err
	}

	if isAvailable {
		dailyBonus, err := mq.SendBonusToBank(bonus)
		if err != nil {
			return dailyBonus, err
		}

		if !conf.IsLocal() {
			err = mqsurvival.Publish(&mqsurvival.TicketGift{Tickets: 3, ToUserId: iden.UserId})
			if err != nil {
				log.Errorf("Failed to publish ticket gift: %s", err)
			}
		}

		return bonus, nil
	} else {
		return nil, nil
	}
}
