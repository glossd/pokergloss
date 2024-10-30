package web

import (
	"fmt"
	"github.com/glossd/pokergloss/bank/domain"
	"github.com/glossd/pokergloss/gomq/mqbank"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToOperation(r *mqbank.DepositRequest) (*domain.Operation, error) {
	var reason string
	if r.Reason != "" {
		reason = r.Reason
	} else {
		depositType, err := mapDepositType(r.Type)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		reason = string(depositType)
	}

	switch r.CurrencyType {
	case mqbank.CurrencyType_COINS:
		operation, err := domain.NewDepositCoins(domain.Reason(reason), r.Chips, r.UserId, r.Description)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return operation, nil
	default:
		operation, err := domain.NewDeposit(domain.Reason(reason), r.Chips, r.UserId, r.Description)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return operation, nil
	}
}

func mapDepositType(depositType mqbank.DepositRequest_Type) (domain.Reason, error) {
	switch depositType {
	case mqbank.DepositRequest_BONUS:
		return domain.Bonus, nil
	case mqbank.DepositRequest_CASH_GAME:
		return domain.Game, nil
	case mqbank.DepositRequest_NEW_LEVEL:
		return domain.NewLevel, nil
	case mqbank.DepositRequest_NEW_ACHIEVEMENT:
		return domain.Achievement, nil
	case mqbank.DepositRequest_ASSIGNMENT:
		return domain.Assignment, nil
	case mqbank.DepositRequest_SURVIVAL:
		return domain.Survival, nil
	}
	log.Errorf("Deposit request type %s is not supported", depositType)
	return "", fmt.Errorf("deposit type %s is not supported", depositType)
}
