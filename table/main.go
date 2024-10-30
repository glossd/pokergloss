package table

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/table/cmd/scheduler"
	"github.com/glossd/pokergloss/table/cmd/server"
	"github.com/glossd/pokergloss/table/db"
	log "github.com/sirupsen/logrus"
	"os"
)

// @title Table API
// @schemes https
// @license.name Pokerblow
// @host pokerblow.com
// @BasePath /api/table
func Run(c *gin.Engine) func(context.Context) {
	appType := os.Getenv("PG_TABLE_SERVICE_TABLE_APP_TYPE")
	switch appType {
	case "scheduler":
		scheduler.Execute()
	case "server":
		server.Execute(c)
	case "all":
		go scheduler.Execute()
		server.Execute(c)
	default:
		log.Warnf("No PG_TABLE_SERVICE_TABLE_APP_TYPE specified, defaulting to all")
		go scheduler.Execute()
		server.Execute(c)
	}
	return func(ctx context.Context) {
		err := db.Client.Disconnect(ctx)
		if err != nil {
			log.Errorf("table db disconnect error: %s", err)
		}
	}
}
