package main

import (
	"flag"
	"github.com/glossd/pokergloss/bank/db"
	"log"
)

func main() {
	userIdP := flag.String("uid", "", "user ID")
	flag.Parse()
	userID := *userIdP
	if userID == "" {
		log.Fatalf("uid can't be empty")
	}

	ctx, dbClient, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Disconnect(ctx)

	db.ReCalcBalance(userID)
}
