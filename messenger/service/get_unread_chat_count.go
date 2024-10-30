package service

import (
	"context"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/messenger/db"
	"github.com/glossd/pokergloss/messenger/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUnreadChatsCount(ctx context.Context, iden authid.Identity) (int64, error) {
	ucl, err := db.FindUserChatList(ctx, iden.UserId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}

	var count int64
	for _, chat := range ucl.U2UChats {
		msg, err := db.FindLastMessageWithChatID(ctx, chat.ID)
		if err != nil {
			log.Errorf("Failed to find chat %s of %s : %s", chat.ID.Hex(), iden, err)
			continue
		}
		if msg.UserID != iden.UserId && msg.Status == domain.Sent {
			count++
		}
	}
	return count, nil
}
