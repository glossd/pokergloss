package bonus

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/bonus/db"
	"github.com/glossd/pokergloss/bonus/services/schedule"
	"github.com/glossd/pokergloss/bonus/web/router"
	log "github.com/sirupsen/logrus"
)

// @title Bonus API
// @schemes https
// @license.name Pokerblow
// @host pokerblow.com
// @BasePath /api/bonus
func Run(c *gin.Engine) func(context.Context) {
	auth.Init()
	_, dbClient, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = schedule.RunBonusesCron()
	if err != nil {
		log.Fatalf("Couldn't start bonus server, %s", err)
	}

	router.New(c)
	return func(ctx context.Context) {
		err := dbClient.Disconnect(ctx)
		if err != nil {
			log.Errorf("db disconnect error: %s", err)
		}
	}
}
