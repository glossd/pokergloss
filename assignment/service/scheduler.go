package service

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"time"
)

func StartDailyScheduler() {
	c := cron.New()
	_, err := c.AddFunc("50 23 * * *", func() {
		tomorrow := time.Now().Add(time.Hour)
		CreateDailyRecursive(tomorrow)
	})
	if err != nil {
		log.Fatalf("Failed to schedule daily: %s", err)
	}
	c.Start()
}

func CreateDailyRecursive(t time.Time) {
	d, err := CreateDaily(t)
	if err != nil {
		log.Errorf("Failed to schedule daily: %s", err)
		time.AfterFunc(10*time.Second, func() {CreateDailyRecursive(t)})
	} else {
		log.Infof("Successfully created new daily %s", d.Day)
	}
}
