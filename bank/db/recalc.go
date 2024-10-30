package db

import (
	"github.com/glossd/pokergloss/bank/domain"
	log "github.com/sirupsen/logrus"
)

// DO NOT USE it in service.
func ReCalcBalance(userID string) {
	ops, err := FindAllOperationsByUserId(userID)
	if err != nil {
		log.Fatalf("Couldn't find user operations: %s", err)
	}

	var chips int64
	var coins int64
	for _, op := range ops {
		switch op.Type {
		case domain.Withdraw:
			chips -= op.Chips
		case domain.Deposit:
			chips += op.Chips
		case domain.DepositCoins:
			coins += op.Coins
		case domain.WithdrawCoins:
			coins -= op.Coins
		}
		if op.Type == domain.Withdraw {
			chips -= op.Chips
		}
		if op.Type == domain.Deposit {
			chips += op.Chips
		}

	}

	b, err := FindBalanceNoCtx(userID)
	if err != nil {
		log.Fatalf("Couldn't find user balance: %s", err)
	}

	b.SetBalanceCalc(chips, coins, ops[len(ops)-1].ID)

	err = UpdateBalanceNoCtx(b)
	if err != nil {
		log.Fatalf("Couldn't upsert user balance: %s", err)
	}

	log.Println("Successfully recalculated")
}
