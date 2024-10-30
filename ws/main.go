package ws

import (
	"context"
	"github.com/gin-gonic/gin"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/ws/db"
	"github.com/glossd/pokergloss/ws/web/fcm"
	"github.com/glossd/pokergloss/ws/web/mq/mqsub"
	"github.com/glossd/pokergloss/ws/web/router"
	log "github.com/sirupsen/logrus"
)

// @title WS service API
// @schemes https
// @license.name Pokerblow
// @host pokerblow.com
// @BasePath /api/ws
func Run(c *gin.Engine) func(context.Context) {
	go mqsub.SubscribeNews()

	_, dbClient, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}

	if conf.IsProd() {
		//meta, err := gke.InitMetadata()
		//if err != nil {
		//	log.Errorf("Failed to init gke.Metadata: %s", err)
		//} else {
		//	metric.PeriodicallyWriteMetrics(meta, map[string]func() int64{
		//		"ws/user_connections": storage.GetUserConnections,
		//	})
		//}
		fcm.Init()
	}

	router.New(c)
	return func(ctx context.Context) {
		err := dbClient.Disconnect(ctx)
		if err != nil {
			log.Errorf("profile db disconnect error: %s", err)
		}
	}
}
