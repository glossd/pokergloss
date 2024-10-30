package tablehistory

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/table-history/db"
	"github.com/glossd/pokergloss/table-history/services/cleaner"
	"github.com/glossd/pokergloss/table-history/web/mq"
	log "github.com/sirupsen/logrus"
)

func Run(c *gin.Engine) func(context.Context) {
	_, dbClient, err := db.Init()
	if err != nil {
		log.Fatalf("Failed to connect to db: %s", err)
	}

	go cleaner.Run()

	go mq.Subscribe()
	return func(ctx context.Context) {
		err := dbClient.Disconnect(ctx)
		if err != nil {
			log.Errorf("db disconnect error: %s", err)
		}
	}
}
