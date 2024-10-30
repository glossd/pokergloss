package scheduler

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/services/cash"
	"github.com/glossd/pokergloss/table/services/cleaning"
	"github.com/glossd/pokergloss/table/services/multi"
	"github.com/glossd/pokergloss/table/services/recovery"
	"github.com/glossd/pokergloss/table/services/sitngo"
	"github.com/glossd/pokergloss/table/web/client/mqpub"
	"github.com/glossd/pokergloss/table/web/client/mqsub"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"time"
)

func Execute() {
	db.Init()

	now := time.Now()
	if now.Hour() > 0 && now.Minute() > 5 { // scheduling start at 00:05
		go recovery.InitMultiSchedulerRecovery()
	}

	if conf.IsProd() || conf.IsLocalOnly() {
		go mqsub.SubscribeForStartMulti()
		go mqsub.SubscribeForStartSitngo()
		go mqsub.SubscribeForMultiRebalance()

		cash.CreatePersistentTables()
		sitngo.CreateDaily()
	}

	c := cron.New()

	_, err := c.AddFunc(conf.Props.TableService.Scheduler.StartMultiTournaments, func() {
		startAt := timeutil.Time(time.Now().Truncate(time.Minute))
		_ = multi.LaunchMultiLobbies(startAt)
	})
	if err != nil {
		log.Fatalf("scheduler: Failed to add start multi lobbies job: %s", err)
	}

	_, err = c.AddFunc("* * * * *", func() {
		startAt := timeutil.Time(time.Now().Truncate(time.Minute))
		mqpub.PublishStartSitngo(startAt)
	})
	if err != nil {
		log.Fatalf("scheduler: Failed to add start sitngo lobbies job: %s", err)
	}

	_, err = c.AddFunc("5 0 * * *", func() {
		tomorrow := time.Now().AddDate(0, 0, 1)
		multi.CreateDailyMultiTournaments(tomorrow, multi.NothingEnrich)
	})
	if err != nil {
		log.Fatalf("scheduler: Failed to add create daily freerolls job: %s", err)
	}

	_, err = c.AddFunc(conf.Props.TableService.Cron.CreateDailySitngo, func() {
		sitngo.CreateDaily()
	})
	if err != nil {
		log.Fatalf("scheduler: Failed to add create daily freerolls job: %s", err)
	}

	_, err = c.AddFunc("30 21 * * *", func() {
		tomorrow := time.Now().AddDate(0, 0, 1)
		multi.CreateTheTournament(tomorrow, multi.NothingEnrich)
	})
	if err != nil {
		log.Fatalf("scheduler: Failed to add create daily freerolls job: %s", err)
	}

	_, err = c.AddFunc(conf.Props.TableService.Cron.CleanSittingOutPlayers, cleaning.CleanSittingOutPlayers)
	if err != nil {
		log.Fatalf("scheduler: Failed to add clean sitting out players job: %s", err)
	}

	_, err = c.AddFunc(conf.Props.TableService.Cron.DeleteFinishedLobbies, cleaning.CleanFinishedLobbies)
	if err != nil {
		log.Fatalf("scheduler: Failed to add delete finished lobbies job: %s", err)
	}

	_, err = c.AddFunc(conf.Props.TableService.Cron.DeleteNotStartedSitngo, cleaning.CleanNotStartedLobbies)
	if err != nil {
		log.Fatalf("scheduler: Failed to add delete not started sitngo job: %s", err)
	}

	_, err = c.AddFunc(conf.Props.TableService.Cron.CleanWaitingTables, func() {
		cleaning.CleanWaitingTables()
	})
	if err != nil {
		log.Fatalf("scheduler: Failed to add clean waiting tables: %s", err)
	}

	_, err = c.AddFunc(conf.Props.TableService.Cron.CleanAlonePlayerOnPersistentTable, func() {
		cleaning.CleanAlonePlayerOnPersistentTable()
	})
	if err != nil {
		log.Fatalf("scheduler: Failed to add clean alone player on a persistent table: %s", err)
	}

	c.Run()
}
