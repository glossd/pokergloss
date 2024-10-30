package idp

import (
	"context"
	"github.com/glossd/pokergloss/gomq"
	log "github.com/sirupsen/logrus"
)

var ErrEmailNotVerified = gomq.NewAckableError("user email is no verified")

func GetUserEmail(ctx context.Context, userID string) (email string, err error) {
	u, err := client.GetUser(ctx, userID)
	if err != nil {
		log.Errorf("idp.GetUserEmail failed: %s", err)
		return "", gomq.WrapInAckableError(err)
	}

	if !u.EmailVerified {
		return "", ErrEmailNotVerified
	}

	return u.Email, nil
}
