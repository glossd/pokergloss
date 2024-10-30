package survival

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/survival/db"
	"github.com/glossd/pokergloss/survival/web/mq"
	"github.com/glossd/pokergloss/survival/web/router"
	log "github.com/sirupsen/logrus"
)

// @title Survival API
// @schemes https
// @license.name Pokerblow
// @host pokerblow.com
// @BasePath /api/survival
func Run(c *gin.Engine) func(ctx context.Context) {
	_, dbClient, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}

	go mq.SubscribeSurvivalEnd()
	go mq.SubscribeTournamentEnd()
	go mq.SubscribeCreatedProfiles()
	go mq.SubscribeTicketGifts()

	router.New(c)
	return func(ctx context.Context) {
		err := dbClient.Disconnect(ctx)
		if err != nil {
			log.Errorf("profile db disconnect error: %s", err)
		}
	}
}
