package bank

import (
	"github.com/glossd/pokergloss/assignment/domain"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq/mqbank"
	log "github.com/sirupsen/logrus"
)

func DepositPrize(userID string, a *domain.Assignment) {
	if !conf.IsProd() {
		return
	}
	err := mqbank.Deposit(&mqbank.DepositRequest{
		Chips:       a.GetPrize(),
		Type:        mqbank.DepositRequest_ASSIGNMENT,
		Description: a.GetFullName(),
		UserId:      userID,
	})
	if err != nil {
		log.Errorf("Failed to send deposit for completed assignment userID=%s, assignmentType=%s: %s", userID, a.Type, err)
	}
}
