package recovery

import (
	"context"
	"github.com/glossd/pokergloss/goconf/timeutil"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/services/multi"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
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

	_, err = db.FindLobbyMulti(context.Background(), multi.SaturdayTournamentID)
	if err == mongo.ErrNoDocuments {
		multi.CreateSaturdayTournament()
	}
}
