package service

import (
	"context"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/messenger/db"
	"github.com/glossd/pokergloss/messenger/domain"
	"github.com/glossd/pokergloss/messenger/web/ws/wssend"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleTyping(ctx context.Context, iden authid.Identity, chatID string, text string) error {
	chatOID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return err
	}
	chat, err := db.FindChat(ctx, chatOID)
	if err != nil {
		log.Errorf("Failed to handle user typing: %s", err)
		return err
	}

	if !chat.ContainsParticipant(iden.UserId) {
		return domain.ErrChatNotAvailable
	}

	for _, participant := range chat.ParticipantsExcept(iden.UserId) {
		wssend.SendTypingTo(participant, chatID, iden)
	}
	return nil
}
