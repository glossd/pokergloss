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

func ReadMessages(ctx context.Context, iden authid.Identity, chatID string, msgIDs []string) error {
	chatOID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return err
	}

	chat, err := db.FindChat(ctx, chatOID)
	if err != nil {
		return err
	}

	if !chat.ContainsParticipant(iden.UserId) {
		return domain.ErrChatNotAvailable
	}

	// todo improve performance
	for _, msgID := range msgIDs {
		msgOID, err := primitive.ObjectIDFromHex(msgID)
		if err != nil {
			return err
		}

		msg, err := db.FindMessage(ctx, msgOID)
		if err != nil {
			return err
		}
		if msg.UserID == iden.UserId {
			return domain.ErrYouCantChangeStatus
		}
		if msg.Status == domain.Read {
			return domain.ErrYouCantChangeStatus
		}
		err = db.UpdateMessageStatus(ctx, msgOID, domain.Read)
		if err != nil {
			return err
		}

		err = wssend.SendStatusTo(msg.UserID, msgID, msg.ChatID.Hex(), domain.Read)
		if err != nil {
			log.Warn("Failed to sent status read")
		}
	}
	return nil
}
