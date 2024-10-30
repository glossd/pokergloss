package mq

import (
	"context"
	"github.com/glossd/pokergloss/gomq/mqprofile"
	"github.com/glossd/pokergloss/gomq/mqsurvival"
	"github.com/glossd/pokergloss/gomq/mqtable"
	"github.com/glossd/pokergloss/survival/db"
	"github.com/glossd/pokergloss/survival/service"
	log "github.com/sirupsen/logrus"
)

func SubscribeSurvivalEnd() {
	err := mqtable.SubscribeSurvivalEnd("survival-service", func(ctx context.Context, msg *mqtable.SurvivalEnd) error {
		return service.EndLevel(ctx, msg.UserId, msg.TableId, msg.IsUserLost)
	})
	if err != nil {
		log.Fatalf("Failed to init survival end subscriber: %s", err)
	}
}

func SubscribeTournamentEnd() {
	err := mqtable.SubscribeTournamentEnd("survival-service-tournament-end", func(ctx context.Context, msg *mqtable.TournamentEnd) error {
		if len(msg.TournamentWinners) == 1 {
			return db.CardIncTicket(ctx, msg.TournamentWinners[0].UserId)
		}

		for _, tw := range msg.TournamentWinners {
			_ = db.CardIncTicket(ctx, tw.UserId)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to init tournament end subscriber: %s", err)
	}
}

func SubscribeTicketGifts() {
	err := mqsurvival.SubscribeForTicketGifts("survival-service-ticket-gifts", func(ctx context.Context, msg *mqsurvival.TicketGift) error {
		if msg.Tickets > 0 {
			return db.GiftTickets(ctx, msg.ToUserId, msg.Tickets)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to init ticket gifts subscriber: %s", err)
	}
}

func SubscribeCreatedProfiles() {
	err := mqprofile.SubscribeForCreated("survival-service-created-profiles", func(ctx context.Context, msg *mqprofile.Profile) error {
		return db.CardIncTwoTickets(ctx, msg.UserId)
	})
	if err != nil {
		log.Fatalf("Failed to init created profiles subscriber: %s", err)
	}
}
