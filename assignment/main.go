package assignment

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/assignment/db"
	"github.com/glossd/pokergloss/assignment/service"
	"github.com/glossd/pokergloss/assignment/web/rest/router"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq/mqtable"
	log "github.com/sirupsen/logrus"
)

// @title Assignment API
// @schemes https
// @license.name Pokerblow
// @host pokerblow.com
// @BasePath /api/assignment
func Run(c *gin.Engine) func(context.Context) {
	_, dbClient, err := db.Init()
	if err != nil {
		log.Fatalf("Failed to connect to db: %s", err)
	}

	service.RecoverDaily()
	service.StartDailyScheduler()

	if conf.IsProd() {
		go func() {
			err = mqtable.SubscribeGameEnd("assignment-service", func(ctx context.Context, end *mqtable.GameEnd) error {
				return service.UpdateDailies(ctx, end)
			})
			if err != nil {
				log.Panicf("Failed to mqtable.SubscribeGameEnd: %s", err)
			}
		}()

		go func() {
			err = mqtable.SubscribeTournamentEnd("assignment-service-tournament-end", func(ctx context.Context, end *mqtable.TournamentEnd) error {
				return service.UpdateDailiesTournament(ctx, end)
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
			log.Errorf("assignment db disconnect errro: %s", err)
		}
	}
}
