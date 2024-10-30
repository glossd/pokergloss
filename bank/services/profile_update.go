package services

import (
	"context"
	"github.com/glossd/pokergloss/bank/db"
	"github.com/glossd/pokergloss/bank/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateProfileInfo(ctx context.Context, info domain.ProfileInfo) error {
	b, err := db.FindBalanceNoCtx(info.UserID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err := CreateFirstBalance(ctx, info)
			if err != nil {
				return err
			}
			return applyWelcomeBonus(ctx, info.UserID)
		}
		log.Errorf("Failed to update profile info, find: %s", err)
		return err
	}
	b.UpdateProfile(info.Username, info.Picture)
	err = db.UpdateProfileInfo(ctx, b)
	if err != nil {
		log.Errorf("Failed to update profile info, update: %s", err)
		return err
	}
	return nil
}

func CreateFirstBalance(ctx context.Context, info domain.ProfileInfo) error {
	err := db.InsertBalance(ctx, domain.NewBalance(info))
	if err != nil {
		log.Errorf("Failed to create first balance: %s", err)
		return err
	}
	return nil
}

func applyWelcomeBonus(ctx context.Context, userID string) error {
	deposit, err := domain.NewDeposit(domain.Bonus, 2500, userID, "Welcome bonus")
	if err != nil {
		return err
	}
	err = Deposit(ctx, deposit)
	if err != nil {
		log.Errorf("Failed to save welcome bonus for userID=%s : %s", userID, err)
		return err
	}
	return nil
}
