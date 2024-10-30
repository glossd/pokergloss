package bank

import (
	"fmt"
	"github.com/glossd/pokergloss/achievement/domain"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq/mqbank"
	log "github.com/sirupsen/logrus"
)

func DepositNewLevel(exp *domain.ExP) {
	if !conf.IsProd() {
		return
	}
	if exp.IsNewLevel() {
		err := mqbank.Deposit(&mqbank.DepositRequest{
			Chips:       exp.GetNewLevelPrize(),
			Type:        mqbank.DepositRequest_NEW_LEVEL,
			Description: fmt.Sprintf("Reached %d level", exp.Level),
			UserId:      exp.UserID,
		})
		if err != nil {
			log.Errorf("Failed to send deposit for new level exp=%+v : %s", exp, err)
		}
	}
}

func DepositHandAchievement(as *domain.AchievementStore) {
	if !conf.IsProd() {
		return
	}
	hc := as.HandsCounter
	if hc.GetPrize().Chips > 0 {
		err := mqbank.Deposit(&mqbank.DepositRequest{
			Type:        mqbank.DepositRequest_NEW_ACHIEVEMENT,
			UserId:      as.UserID,
			Chips:       hc.GetPrize().Chips,
			Description: fmt.Sprintf("Achievement: won with %s %d times", hc.GetPrize().Name, hc.GetPrizeHandCount()),
		})
		if err != nil {
			log.Errorf("Failed to send prize userId %s, prize %+v : %s", as.UserID, hc.GetPrize(), err)
		}
	}
}

func DepositCounterAchievement(userID string, c domain.Counter) {
	if !conf.IsProd() {
		return
	}
	if c.GetPrize().Chips == 0 {
		return
	}
	err := mqbank.Deposit(&mqbank.DepositRequest{
		Type:        mqbank.DepositRequest_NEW_ACHIEVEMENT,
		UserId:      userID,
		Chips:       c.GetPrize().Chips,
		Description: fmt.Sprintf("Achievement: %s %d times", c.GetName(), c.GetCount()),
	})
	if err != nil {
		log.Errorf("Failed to send prize userId %s, prize %+v : %s", userID, c.GetPrize(), err)
	}
}
