package cleaner

import (
	"github.com/glossd/pokergloss/table-history/db"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func Run() {
	c := cron.New()
	_, err := c.AddFunc("0 2 * * *", func() {
		db.DeleteAllBefore(time.Now().AddDate(0, 0, -7))
	})
	if err != nil {
		log.Fatalf("Cron expression is wrong: %s", err)
	}
	c.Run()
}
