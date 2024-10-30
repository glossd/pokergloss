package service

import (
	"context"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/messenger/db"
	"github.com/glossd/pokergloss/messenger/domain"
	"github.com/glossd/pokergloss/messenger/web/clients/profile"
	"github.com/glossd/pokergloss/messenger/web/model"
	log "github.com/sirupsen/logrus"
	"time"
)

func GetUserChats(ctx context.Context, iden authid.Identity) ([]*model.Chat, error) {
	ulc, err := findOrBuildUserChatList(ctx, iden.UserId)
	if err != nil {
		return nil, err
	}

	for _, chat := range ulc.U2UChats {
		lastMessage, err := db.FindLastMessageWithChatID(ctx, chat.ID)
		if err != nil {
			log.Errorf("Failed to find last chat message of chat %s :%s", chat.ID, err)
			continue
		}
		chat.SetLastMessage(lastMessage)
	}

	userMap := profile.FetchUserMap(ctx, ulc)

	var chats []*model.Chat
	if !ulc.IsSupportCreated {
		userID := domain.SupportUserID
		status := domain.Read
		text := "User Support. If you have any questions, please ask them here."
		updatedAt := time.Now().UnixNano() / 1e6
		chatID := "phantom-support"
		chats = append(chats, &model.Chat{
			ID:      chatID,
			Name:    "support",
			Picture: "https://storage.googleapis.com/avatarsforpoker/8xJcx5LfCrTXUsaRdHC8ufzmjdT2-uP30",
			LastMessage: &model.Message{
				ChatID:    &chatID,
				UserID:    &userID,
				Status:    &status,
				Text:      &text,
				UpdatedAt: &updatedAt,
			},
			IsPhantom: true,
		})
	}
	for _, chat := range ulc.GetSortedChats() {
		identity, ok := userMap[chat.OtherUserID]
		if ok {
			c := &model.Chat{
				ID:          chat.ID.Hex(),
				LastMessage: model.ToMessage(chat.GetLastMessage()),
				Name:        identity.Username,
				Picture:     identity.Picture,
			}
			chats = append(chats, c)
		}
	}
	return chats, nil
}
