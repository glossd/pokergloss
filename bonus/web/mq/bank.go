package mq

import (
	"github.com/glossd/pokergloss/bonus/domain"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq/mqbank"
)

func SendBonusToBank(bonus *domain.DailyBonus) (*domain.DailyBonus, error) {
	if !conf.IsLocal() {
		r := mqbank.DepositRequest{
			Chips:       bonus.CalculateBonus(),
			Type:        mqbank.DepositRequest_BONUS,
			Description: "Daily bonus",
			UserId:      bonus.UserID,
		}
		err := mqbank.Deposit(&r)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
