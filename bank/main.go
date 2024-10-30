package bank

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/bank/db"
	"github.com/glossd/pokergloss/bank/services/fee"
	"github.com/glossd/pokergloss/bank/services/ranker"
	"github.com/glossd/pokergloss/bank/web/mq/mqsub"
	"github.com/glossd/pokergloss/bank/web/router"
	conf "github.com/glossd/pokergloss/goconf"
	log "github.com/sirupsen/logrus"
	"time"
)

// @title Bank API
// @schemes https
// @license.name Pokerblow
// @host pokerblow.com
// @BasePath /api/bank
func Run(c *gin.Engine) func(context.Context) {
	_, dbClient, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}

	go mqsub.SubscribeToDeposits()
	go mqsub.SubscribeToProfiles()
	go mqsub.SubscribeForBalanceUpdates()

	ticker := time.NewTicker(conf.Props.RankerDuration)
	go ranker.RunRanker(ticker)
	fee.LaunchInactionFeeJob()

	router.New(c)
	return func(ctx context.Context) {
		err := dbClient.Disconnect(ctx)
		if err != nil {
			log.Errorf("bank failed to disconnect db: %s", err)
		}
	}
}
