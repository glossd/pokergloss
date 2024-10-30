package job

import (
	"context"
	"github.com/glossd/pokergloss/achievement/db"
	"github.com/glossd/pokergloss/achievement/domain"
	"github.com/glossd/pokergloss/gomq/mqtable"
	log "github.com/sirupsen/logrus"
)

type Lock struct {
	ID string `bson:"_id"`
}

func Run() {
	log.Info("Starting the job")
	_, err := db.Client.Database(db.DbName).Collection("jobBust").InsertOne(context.Background(), &Lock{ID: "lock"})
	if err != nil {
		log.Errorf("Failed to insert lock")
		return
	}
	err = db.ForeachGameEnd(func(end *mqtable.GameEnd) {
		ge := domain.NewGameEnd(end)
		for _, winner := range end.Winners {
			as, err := db.FindAchievementStoreNoCtx(winner.UserId)
			if err != nil {
				continue
			}
			if as.DefeatCounter == nil {
				as.DefeatCounter = domain.NewDefeatCounter()
			}
			as.DefeatCounter.Update(winner, ge)
			if as.BustCounter == nil {
				as.BustCounter = domain.NewBustCounter()
			}
			as.BustCounter.Update(winner, ge)
			err = db.UpsertAchievementStore(context.Background(), as)
			if err != nil {
				log.Errorf("Failed to upsert achievement store of userId=%s : %s", winner.UserId, err)
			}
		}
	})
	if err != nil {
		log.Errorf("Job failed: %s", err)
	} else {
		log.Info("Successfully finished the job")
	}
}
