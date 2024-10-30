package service

import (
	"context"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/messenger/db"
	"github.com/glossd/pokergloss/messenger/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetChatMessages(ctx context.Context, iden authid.Identity, chatID string, lastID string, limit int64) ([]*domain.Message, error) {
	if chatID == "phantom-support" {
		return []*domain.Message{}, nil
	}
	chatOID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return nil, err
	}
	var lastOID primitive.ObjectID
	if lastID != "" {
		lastOID, err = primitive.ObjectIDFromHex(lastID)
		if err != nil {
			return nil, err
		}
	}
	chat, err := db.FindChat(ctx, chatOID)
	if err != nil {
		return nil, err
	}
	if !chat.ContainsParticipant(iden.UserId) {
		return nil, domain.ErrChatNotAvailable
	}
	return db.FindMessagesByChatID(ctx, chatOID, lastOID, limit)
}
