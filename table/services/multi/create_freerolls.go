package multi

import (
	"fmt"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const NumberOfTournaments = 3 + 1

var NothingEnrich = func(multi *domain.LobbyMulti) {}

func CreateDailyMultiTournaments(t time.Time, enricher func(multi *domain.LobbyMulti), except ...time.Time) {
	midnight := t.In(time.UTC).Truncate(24 * time.Hour)
	createTournaments(t)

	db.UpdateLastCreatedHourlyFreerolls(midnight)
	log.Infof("Successfully created daily freerolls")
}

func CreateTheTournament(t time.Time, enricher func(multi *domain.LobbyMulti), except ...time.Time) {
	createFreerolls(t, enricher, except)
}

func createTournaments(t time.Time) {
	midnight := t.In(time.UTC).Truncate(24 * time.Hour)

	freerollTimes, err := timeutil.TimeSeriesOfDay(t, conf.Props.Multi.Freerolls.At...)
	if err != nil {
		log.Panicf("Failed to create freerolls, time in config is wrong: %s", err)
	}

	var lobbies []interface{}

	// we don't have many users, just 3 tournaments a day is enough.
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("Tournament %d #%d", midnight.Day(), i+1)
		startAt := midnight.Add(time.Duration(i)*8*time.Hour + 4*time.Hour + 15*time.Minute) // 04:15, 12:15, 20:15
		if timeutil.Contains(freerollTimes, startAt) {
			continue
		}
		lobbies = append(lobbies, domain.NewMultiLobbyDaily(name, domain.NL, startAt, 1000, 10))
	}

	//// cheap
	//for i := 0; i < 24; i++ {
	//	name := fmt.Sprintf("Tournament %d #%d", midnight.Day(), i+1)
	//	startAt := midnight.Add(time.Duration(i) * time.Hour)
	//	if timeutil.Contains(freerollTimes, startAt) {
	//		continue
	//	}
	//	lobbies = append(lobbies, domain.NewMultiLobbyDaily(name, domain.NL, startAt, 1000, 10))
	//}

	//// medium
	//for i := 0; i < 8; i++ {
	//	name := fmt.Sprintf("Medium %d #%d", midnight.Day(), i+1)
	//	startAt := midnight.Add(time.Duration(i)*2*time.Hour + time.Hour + time.Hour/2) // 01:30 - 22:30
	//	lobbies = append(lobbies, domain.NewMultiLobbyDaily(name, domain.NL, startAt, 5000, 40))
	//}
	//
	//// rich
	//for i := 0; i < 3; i++ {
	//	name := fmt.Sprintf("High %d #%d", midnight.Day(), i+1)
	//	startAt := midnight.Add(time.Duration(i)*8*time.Hour + 4*time.Hour + 15*time.Minute) // 04:15, 12:15, 20:15
	//	lobbies = append(lobbies, domain.NewMultiLobbyDaily(name, domain.ML, startAt, 25000, 100))
	//}
	//
	//// super one
	//name := fmt.Sprintf("Super %d", midnight.Day())
	//startAt := midnight.Add(21*time.Hour + 45*time.Minute) // 21:45
	//lobbies = append(lobbies, domain.NewMultiLobbyDaily(name, domain.ML, startAt, 100000, 200))

	err = db.InsertManyLobbyMulti(lobbies)
	if err != nil {
		log.Errorf("Failed insert multi lobbies: %s", err)
	}
}

func createFreerolls(t time.Time, enricher func(multi *domain.LobbyMulti), except []time.Time) {
	var lobbies []interface{}
	freerollTimes, err := timeutil.TimeSeriesOfDay(t, conf.Props.Multi.Freerolls.At...)
	if err != nil {
		log.Panicf("Failed to create freerolls, time in config is wrong: %s", err)
	}
	for i, startAt := range freerollTimes {
		if startAt.Before(time.Now()) {
			continue
		}
		if timeutil.Contains(except, startAt) {
			continue
		}

		freeroll := domain.NewFreerollWithNumber(startAt, i+1)
		enricher(freeroll)
		if startAt.Hour() == 19 && startAt.Minute() == 0 {
			oid, err := primitive.ObjectIDFromHex("60280c0072d38a27d9264d88")
			if err != nil {
				log.Fatal(err)
			}
			freeroll.ID = oid
		}
		lobbies = append(lobbies, freeroll)
	}

	err = db.InsertManyLobbyMulti(lobbies)
	if err != nil {
		log.Errorf("Failed insert multi lobbies: %s", err)
	}
}
