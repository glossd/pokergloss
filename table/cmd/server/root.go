package server

import (
	"github.com/gin-gonic/gin"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/web/client/mqpub"
	"github.com/glossd/pokergloss/table/web/client/mqsub"
	"github.com/glossd/pokergloss/table/web/router"
	"math/rand"
	"time"
)

func Execute(c *gin.Engine) {
	// adds randomness to domain.Algo https://stackoverflow.com/a/12321192/10160865
	rand.Seed(time.Now().UTC().UnixNano())
	db.Init()
	if conf.IsProd() || conf.IsLocalOnly() {
		// uncomment when you have any metrics
		//err := metrics.Init()
		//if err != nil {
		//	log.Fatalf("Failed connection to StackDriver server: %s", err)
		//}

		mqpub.InitTimeoutPublisher()

		go mqsub.SubscribeForTimeouts()
	}

	router.New(c)
}
