package main

import (
	"github.com/glossd/pokergloss/bank/db"
	log "github.com/sirupsen/logrus"
)

func main() {

	ctx, dbClient, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Disconnect(ctx)

	//fee.RunInactionFee()
}
