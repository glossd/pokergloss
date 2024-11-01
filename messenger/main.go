package messenger

import (
	"context"
	"github.com/gin-gonic/gin"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/messenger/db"
	"github.com/glossd/pokergloss/messenger/web/mq/mqsub"
	"github.com/glossd/pokergloss/messenger/web/router"
	log "github.com/sirupsen/logrus"
)

// @title Messenger API
// @schemes https
// @license.name Pokerblow
// @host pokerblow.com
// @BasePath /api/messenger
func Run(c *gin.Engine) func(context.Context) {
	_, dbClient, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}

	if conf.IsProd() {
		go mqsub.Subscribe()
	}

	router.New(c)
	return func(ctx context.Context) {
		err := dbClient.Disconnect(ctx)
		if err != nil {
			log.Errorf("market db disconnect errro: %s", err)
		}
	}
}
