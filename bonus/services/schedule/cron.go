package schedule

import (
	"github.com/glossd/pokergloss/bonus/db"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/robfig/cron/v3"
)

func RunBonusesCron() error {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc(conf.Props.Bonus.Cron, db.ResetBonuses)
	if err != nil {
		return err
	}
	c.Start()
	return nil
}
