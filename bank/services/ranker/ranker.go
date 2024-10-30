package ranker

import (
	"github.com/glossd/pokergloss/bank/db"
	"time"
)

func RunRanker(ticker *time.Ticker) {
	// do-while loop semantic
	run()
	for range ticker.C {
		run()
	}
}

func run() {
	db.BuildOppositeRankView()
}
