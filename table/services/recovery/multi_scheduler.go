package recovery

import (
	"github.com/glossd/pokergloss/goconf/timeutil"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/services/multi"
	log "github.com/sirupsen/logrus"
	"time"
)

func InitMultiSchedulerRecovery() {
	multiConfig, err := db.GetMultiConfig()
	if err != nil {
		log.Errorf("Failed to apply multi table scheduler recovery: %s", err)
		return
	}

	if multiConfig.CreatedHourlyFreerollsAt < timeutil.Midnight(time.Now()) {
		multi.CreateDailyMultiTournaments(time.Now(), multi.NothingEnrich)
	}

	nextDay := time.Now().AddDate(0, 0, 1)
	if multiConfig.CreatedHourlyFreerollsAt < timeutil.Midnight(nextDay) {
		multi.CreateDailyMultiTournaments(nextDay, multi.NothingEnrich)
	}
}
