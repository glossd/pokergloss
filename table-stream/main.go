package tablestream

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table-stream/web/mq/mqsub"
	"github.com/glossd/pokergloss/table-stream/web/router"
)

// @title Table-Stream service API
// @schemes https
// @license.name Pokerblow
// @host pokerblow.com
// @BasePath /api/table-stream
func Run(c *gin.Engine) func(context.Context) {
	auth.Init()
	go mqsub.SubscribeTable()

	if conf.IsProd() {
		//meta, err := gke.InitMetadata()
		//if err != nil {
		//	log.Errorf("Failed to init gke.Metadata: %s", err)
		//} else {
		//	metric.PeriodicallyWriteMetrics(meta, map[string]func() int64{
		//		"ws/table_connections": ws.GetTableConnsCount,
		//	})
		//}
	}

	router.New(c)
	return func(ctx context.Context) {}
}
