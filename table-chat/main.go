package tablechat

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/table-chat/db"
	"github.com/glossd/pokergloss/table-chat/web/rest/router"
	log "github.com/sirupsen/logrus"
)

// @title Table Chat API
// @schemes https
// @license.name Pokerblow
// @host pokerblow.com
// @BasePath /api/table-chat
func Run(c *gin.Engine) func(context.Context) {
	auth.Init()
	_, dbClient, err := db.Init()
	if err != nil {
		log.Fatalf("Failed to connect to db: %s", err)
	}

	router.New(c)
	return func(ctx context.Context) {
		err := dbClient.Disconnect(ctx)
		if err != nil {
			log.Errorf("db disconnect error: %s", err)
		}
	}
}
