package main

import (
	"context"
	"flag"
	"github.com/glossd/pokergloss/bank/db"
	"github.com/glossd/pokergloss/bank/domain"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	userIdP := flag.String("uid", "", "user ID")

	value := flag.Int64("value", 0, "chips or coins")

	opType := flag.String("type", "w", "deposit or withdraw")

	flag.Parse()
	userID := *userIdP
	if userID == "" {
		log.Fatalf("uid can't be empty")
	}
	if *value == 0 {
		log.Fatalf("must set value")
	}

	ctx, dbClient, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Disconnect(ctx)

	var op *domain.Operation
	if *opType == "w" {
		op, err = domain.NewWithdraw(domain.Admin, *value, userID, "Action made by admin")
		if err != nil {
			log.Fatal(err)
		}
	}

	if *opType == "wc" {
		op, err = domain.NewWithdrawCoins(domain.Admin, *value, userID, "Action made by admin")
		if err != nil {
			log.Fatal(err)
		}
	}

	if *opType == "d" {
		op, err = domain.NewDeposit(domain.Admin, *value, userID, "Action made by admin")
		if err != nil {
			log.Fatal(err)
		}
	}

	if *opType == "dc" {
		op, err = domain.NewDepositCoins(domain.Admin, *value, userID, "Action made by admin")
		if err != nil {
			log.Fatal(err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = db.InsertOperation(ctx, op)
	if err != nil {
		log.Fatalf("Couldn't insert operation: %s", err)
	}

	log.Println("Success op")

	db.ReCalcBalance(userID)
}
