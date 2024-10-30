package profile

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/profile/db"
	"github.com/glossd/pokergloss/profile/web/client/authclient"
	"github.com/glossd/pokergloss/profile/web/client/gcs"
	"github.com/glossd/pokergloss/profile/web/router"
	log "github.com/sirupsen/logrus"
)

// @title Profile API
// @schemes https
// @license.name Pokerblow
// @host pokerblow.com
// @BasePath /api/profile
func Run(c *gin.Engine) func(context.Context) {
	_, dbClient, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}

	authclient.Init()

	if goconf.IsProd() {
		gcsCancel := gcs.Init()
		defer gcsCancel()
	}

	router.New(c)
	return func(ctx context.Context) {
		err := dbClient.Disconnect(ctx)
		if err != nil {
			log.Errorf("profile db disconnect error: %s", err)
		}
	}
}
