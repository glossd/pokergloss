package service

import (
	"context"
	"github.com/glossd/pokergloss/profile/db"
	"github.com/glossd/pokergloss/profile/domain"
)

func GetProfile(ctx context.Context, username string) (*domain.Profile, error) {
	return db.FindProfile(ctx, username)
}
