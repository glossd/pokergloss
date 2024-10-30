package achievement

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/achievement/db"
	"github.com/glossd/pokergloss/achievement/service"
	"github.com/glossd/pokergloss/achievement/web/rest/router"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq/mqtable"
	log "github.com/sirupsen/logrus"
)

const SubID = "achievement-service"

// @title Achievement API
// @schemes https
// @license.name PokerGloss
// @host pokergloss.com
// @BasePath /api/achievement
func Run(c *gin.Engine) func(ctx context.Context) {
	_, dbClient, err := db.Init()
	if err != nil {
		log.Fatalf("Failed to connect to db: %s", err)
	}

	if conf.IsProd() {

		go func() {
			err = mqtable.SubscribeGameEnd(SubID, func(ctx context.Context, end *mqtable.GameEnd) error {
				service.UpdateExp(end)
				service.UpdateAchievementStore(ctx, end)
				service.UpdateStatistics(end)
				_ = db.InsertGameEnd(ctx, end)
				return nil
			})
			if err != nil {
				log.Panicf("Failed to mqtable.SubscribeGameEnd: %s", err)
			}
		}()

		go func() {
			err = mqtable.SubscribeTournamentEnd("achievement-service-tournament-end", func(ctx context.Context, end *mqtable.TournamentEnd) error {
				service.UpdateAchievementStoreTournament(ctx, end)
				_ = db.InsertTournamentEnd(ctx, end)
				return nil
			})
			if err != nil {
				log.Panicf("Failed to mqtable.SubscribeTournamentEnd: %s", err)
			}
		}()
	}

	router.New(c)

	return func(ctx context.Context) {
		err := dbClient.Disconnect(ctx)
		if err != nil {
			log.Errorf("achievement db disconnect error: %s", err)
		}
	}
	//go job.Run()
}
