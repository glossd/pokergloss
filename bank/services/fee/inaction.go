package fee

import (
	"context"
	"github.com/glossd/pokergloss/bank/db"
	"github.com/glossd/pokergloss/bank/domain"
	"github.com/glossd/pokergloss/bank/services"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"time"
)

func LaunchInactionFeeJob() {
	c := cron.New()
	_, err := c.AddFunc("0 0 * * *", RunInactionFee)
	if err != nil {
		log.Fatalf("Failed cron inaction fee: %s", err)
	}
	c.Start()
}

func RunInactionFee() {
	ctx := context.Background()
	log.Infof("Starting inaction fee job")
	err := db.ForEachBalance(ctx, func(b *domain.Balance) {
		if b.Chips < 50000 {
			return
		}
		now := time.Now()
		midnight := now.Truncate(24 * time.Hour)
		due := midnight.AddDate(0, 0, -30).UnixNano() / 1e6
		if b.UpdatedAt < due {
			fee := int64(float64(b.Chips) * 0.05)
			withdraw, err := domain.NewWithdraw(domain.InactionFee, fee, b.UserID, "Inaction fee")
			if err != nil {
				log.Errorf("Failed to create withdraw for inaction fee: %s", err)
				return
			}
			err = services.Withdraw(ctx, withdraw)
			if err != nil {
				log.Errorf("Failed to withdraw inaction fee: %s", err)
			}
			log.Infof("Took inaction fee from {%s,%s} %d chips", b.UserID, b.Username, fee)
		}
	})
	if err != nil {
		log.Errorf("RunInactionFee failed: %s", err)
		return
	}
	log.Infof("Successfully finished inaction fee job")
}
