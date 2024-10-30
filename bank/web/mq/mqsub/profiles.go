package mqsub

import (
	"context"
	"github.com/glossd/pokergloss/bank/domain"
	"github.com/glossd/pokergloss/bank/services"
	"github.com/glossd/pokergloss/gomq/mqprofile"
	"log"
)

func SubscribeToProfiles() {
	err := mqprofile.Subscribe("bank-service", func(ctx context.Context, msg *mqprofile.Profile) error {
		return services.UpdateProfileInfo(ctx, domain.ProfileInfo{
			UserID:   msg.UserId,
			Username: msg.Username,
			Picture:  msg.Picture,
		})
	})
	if err != nil {
		log.Fatalf("SubscribeToProfiles failed: %s", err)
	}
}
